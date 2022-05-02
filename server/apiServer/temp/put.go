package temp

import (
	"distributed_storage_system/server/apiServer/locate"
	"distributed_storage_system/utils/elasticSearch"
	"distributed_storage_system/utils/headutils"
	"distributed_storage_system/utils/rs"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

/*
PUT /temp/token
head: range: bytes=<first>-<hash>
body: 文件内容
将body写入token对应的临时文件
*/
func put(w http.ResponseWriter, r *http.Request) {
	//从url获取token
	token := strings.Split(r.URL.EscapedPath(), "/")[2]
	//根据token生成RSResumablePutStream，用于写入数据
	stream, e := rs.NewRSResumablePutStreamFromToken(token)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusForbidden)
		return
	}
	//获取第一个分片的大小
	current := stream.CurrentSize()
	if current == -1 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	//获取偏移量
	offset := headutils.GetOffsetFromHeader(r.Header)
	//比较偏移量和当前文件大小无错误后继续存储
	if current != offset {
		w.WriteHeader(http.StatusRequestedRangeNotSatisfiable)
		return
	}
	//lem(bytes) = 8000*4
	bytes := make([]byte, rs.BLOCK_SIZE)
	for {
		n, e := io.ReadFull(r.Body, bytes)
		if e != nil && e != io.EOF && e != io.ErrUnexpectedEOF {
			log.Println(e)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		current += int64(n)
		//如当前大小超过预设大小,删除临时文件后报错
		if current > stream.Size {
			stream.Commit(false)
			log.Println("resumable put exceed size")
			w.WriteHeader(http.StatusForbidden)
			return
		}
		if n != rs.BLOCK_SIZE && current != stream.Size {
			return
		}
		//将数据写入cache缓存,缓存够大就flush
		stream.Write(bytes[:n])
		//如果当前写入的数据已经等于预计写入数据
		if current == stream.Size {
			//将数据从缓存写入临时文件
			stream.Flush()
			//创建decoder
			//reader：每个server的uuid对应的用于保存数据的临时文件的内容
			getStream, e := rs.NewRSResumableGetStream(stream.Servers, stream.Uuids, stream.Size)
			//校验临时文件hash是否正确
			hash := url.PathEscape(headutils.CalculateHash(getStream))
			//不正确则删除临时文件报错
			if hash != stream.Hash {
				stream.Commit(false)
				log.Println("resumable put done but hash mismatch")
				w.WriteHeader(http.StatusForbidden)
				return
			}
			//判断是否已经保存相同内容文件
			if locate.Exist(url.PathEscape(hash)) {
				stream.Commit(false)
			} else {
				stream.Commit(true)
			}
			//保存元数据
			e = elasticSearch.AddVersion(stream.Name, stream.Hash, stream.Size)
			if e != nil {
				log.Println(e)
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}
	}
}
