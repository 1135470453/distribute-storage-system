package rs

import (
	"distributed_storage_system/utils/objectStream"
	"fmt"
	"io"
	"log"
)

type RSPutStream struct {
	*encoder
}

func NewRSPutStream(dataServers []string, hash string, size int64) (*RSPutStream, error) {
	log.Println("NewRSPutStream start")
	//检查server数量
	if len(dataServers) != ALL_SHARDS {
		return nil, fmt.Errorf("dataServers number mismatch")
	}
	//计算每个分片的大小(size/4向上取整)
	perShard := (size + DATA_SHARDS - 1) / DATA_SHARDS
	writers := make([]io.Writer, ALL_SHARDS)
	var e error
	//将每个分片存到
	for i := range writers {
		//为每个分片创建对应的临时文件
		writers[i], e = objectStream.NewTempPutStream(dataServers[i],
			fmt.Sprintf("%s.%d", hash, i), perShard)
		if e != nil {
			return nil, e
		}
	}
	enc := NewEncoder(writers)
	log.Println("NewRSPutStream end")
	return &RSPutStream{enc}, nil
}

func (s *RSPutStream) Commit(success bool) {
	log.Println("RSPutStream commit start")
	s.Flush()
	for i := range s.writers {
		s.writers[i].(*objectStream.TempPutStream).Commit(success)
	}
	log.Println("RSPutStream commit end")
}
