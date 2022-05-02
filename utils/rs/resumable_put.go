package rs

import (
	"distributed_storage_system/utils/headutils"
	"distributed_storage_system/utils/objectStream"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

//可以与token进行转化,使用该结构可以生成RSResumablePutStream
type resumableToken struct {
	Name    string
	Size    int64
	Hash    string
	Servers []string //可以使用Servers和Uuids生成对应的TempPutStream用于写入数据
	Uuids   []string
}

//用于向临时文件分片写入数据
type RSResumablePutStream struct {
	*RSPutStream    //保存编码器和写入数据的writer
	*resumableToken //保存uuid
}

//返回一个RSResumablePutStream(保存编码器和写入数据、writer、token)
func NewRSResumablePutStream(dataServers []string, name, hash string, size int64) (*RSResumablePutStream, error) {
	putStream, e := NewRSPutStream(dataServers, hash, size)
	if e != nil {
		return nil, e
	}
	//将每个分片的uuid保存在uuids中
	uuids := make([]string, ALL_SHARDS)
	for i := range uuids {
		uuids[i] = putStream.writers[i].(*objectStream.TempPutStream).Uuid
	}

	token := &resumableToken{name, size, hash, dataServers, uuids}
	return &RSResumablePutStream{putStream, token}, nil
}

//根据token生成RSResumablePutStream，用于写入数据
func NewRSResumablePutStreamFromToken(token string) (*RSResumablePutStream, error) {
	//对token进行BASE64解码,并还原为resumableToken t
	b, e := base64.StdEncoding.DecodeString(token)
	if e != nil {
		return nil, e
	}
	var t resumableToken
	e = json.Unmarshal(b, &t)
	if e != nil {
		return nil, e
	}
	//创建为
	writers := make([]io.Writer, ALL_SHARDS)
	for i := range writers {
		writers[i] = &objectStream.TempPutStream{t.Servers[i], t.Uuids[i]}
	}
	enc := NewEncoder(writers)
	return &RSResumablePutStream{&RSPutStream{enc}, &t}, nil
}

//将RSResumablePutStream转化为token
func (s *RSResumablePutStream) ToToken() string {
	b, _ := json.Marshal(s)
	return base64.StdEncoding.EncodeToString(b)
}

//获取当前文件的大小
func (s *RSResumablePutStream) CurrentSize() int64 {
	//获取第一个分片的大小
	r, e := http.Head(fmt.Sprintf("http://%s/temp/%s", s.Servers[0], s.Uuids[0]))
	if e != nil {
		log.Println(e)
		return -1
	}
	if r.StatusCode != http.StatusOK {
		log.Println(r.StatusCode)
		return -1
	}
	size := headutils.GetSizeFromHeader(r.Header) * DATA_SHARDS
	if size > s.Size {
		size = s.Size
	}
	return size
}
