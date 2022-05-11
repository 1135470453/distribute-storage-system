package rs

import (
	"distributed_storage_system/utils/objectStream"
	"fmt"
	"io"
	"log"
)

//用于上传文件转换为RS纠删码分片后存储在对应的dataServer
type RSPutStream struct {
	*encoder
}

//创建一个包含编码器和写入writer的RSPutStream
func NewRSPutStream(dataServers []string, hash string, size int64) (*RSPutStream, error) {
	log.Println("NewRSPutStream start")
	//检查server数量
	if len(dataServers) != ALL_SHARDS {
		return nil, fmt.Errorf("dataServers number mismatch")
	}
	//计算每个分片的大小(size/4向上取整)
	perShard := (size + DATA_SHARDS - 1) / DATA_SHARDS
	//六个用于给分片写入数据
	writers := make([]io.Writer, ALL_SHARDS)
	var e error
	//将每个分片存到
	for i := range writers {
		//让dataServer创建存储文件基本信息的uuid文件,和用于缓存文件内容的uuid.dat文件,并返回TempPutStream格式的uuid和server
		writers[i], e = objectStream.NewTempPutStream(dataServers[i], fmt.Sprintf("%s.%d", hash, i), perShard)
		if e != nil {
			return nil, e
		}
	}
	//创建一个包括编码器和writer的ecoder
	enc := NewEncoder(writers)
	log.Println("NewRSPutStream end")
	return &RSPutStream{enc}, nil
}

//将缓存写入临时文件. true则将临时文件变为正式文件, false则删除临时文件
func (s *RSPutStream) Commit(success bool) {
	log.Println("RSPutStream commit start")
	s.Flush()
	for i := range s.writers {
		s.writers[i].(*objectStream.TempPutStream).Commit(success)
	}
	log.Println("RSPutStream commit end")
}
