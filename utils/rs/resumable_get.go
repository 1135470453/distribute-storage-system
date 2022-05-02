package rs

import (
	"distributed_storage_system/utils/objectStream"
	"io"
)

//用于获取临时文件
type RSResumableGetStream struct {
	*decoder
}

/*
创建decoder
reader：每个server的uuid对应的用于保存数据的临时文件的内容
*/
func NewRSResumableGetStream(dataServers []string, uuids []string, size int64) (*RSResumableGetStream, error) {
	readers := make([]io.Reader, ALL_SHARDS)
	var e error
	//获取每个server的uuid对应的用于保存数据的临时文件的内容
	for i := 0; i < ALL_SHARDS; i++ {
		readers[i], e = objectStream.NewTempGetStream(dataServers[i], uuids[i])
		if e != nil {
			return nil, e
		}
	}
	writers := make([]io.Writer, ALL_SHARDS)
	dec := NewDecoder(readers, writers, size)
	return &RSResumableGetStream{dec}, nil
}
