package rs

import (
	"distributed_storage_system/utils/objectStream"
	"fmt"
	"io"
	"log"
)

type RSGetStream struct {
	*decoder
}

/*
locateInfo:储存文件分片的dataServer地址,[分片号]:dataServer地址
dataServers:用于保存修复节点的dataServer地址
hash:文件对应hash值
size:文件对应的大小
返回RSGetStream,RSGetStream内嵌decoder,decoder中reader保存正确的分片文件,writer用于处理不正确的分片
*/
func NewRSGetStream(locateInfo map[int]string, dataServers []string, hash string, size int64) (*RSGetStream, error) {
	log.Println("NewRSGetStream start")
	if len(locateInfo)+len(dataServers) != ALL_SHARDS {
		return nil, fmt.Errorf("dataServers number mismatch")
	}
	//用于存从dataserver读取的文件
	readers := make([]io.Reader, ALL_SHARDS)
	for i := 0; i < ALL_SHARDS; i++ {
		server := locateInfo[i]
		//存在文件损坏,将新节点位置添加到locateInfo中
		if server == "" {
			locateInfo[i] = dataServers[0]
			dataServers = dataServers[1:]
			continue
		}
		//向Server对应的dataServe发出请求，获取以GetStream形式的已经进行内容检验、被压缩的object
		reader, e := objectStream.NewGetSteam(server, fmt.Sprintf("%s.%d", hash, i))
		//将file内容存入reader
		if e == nil {
			readers[i] = reader
		}
	}
	writers := make([]io.Writer, ALL_SHARDS)
	//获取每个分片的数据大小
	//size/6向上取整
	perShard := (size + DATA_SHARDS - 1) / DATA_SHARDS
	var e error
	for i := range readers {
		//存在损坏文件
		if readers[i] == nil {
			//让dataServer建立临时文件,并返回对应的uuid(uuid为随机生成，用于临时文件名称)
			writers[i], e = objectStream.NewTempPutStream(locateInfo[i], fmt.Sprintf("%s.%d", hash, i), perShard)
			if e != nil {
				return nil, e
			}
		}
	}
	//正确的分片保存在readers中,不正确的分片用writer记录
	dec := NewDecoder(readers, writers, size)
	log.Println("NewRSGetStream end")
	return &RSGetStream{dec}, nil
}

func (s *RSGetStream) Close() {
	log.Println("rs get Close start")
	for i := range s.writers {
		if s.writers[i] != nil {
			s.writers[i].(*objectStream.TempPutStream).Commit(true)
		}
	}
	log.Println("rs get Close end")
}

//读取到偏移量位置
func (s *RSGetStream) Seek(offset int64, whence int) (int64, error) {
	if whence != io.SeekCurrent {
		panic("only support SeekCurrent")
	}
	if offset < 0 {
		panic("only support forward seek")
	}
	for offset != 0 {
		length := int64(BLOCK_SIZE)
		if offset < length {
			length = offset
		}
		buf := make([]byte, length)
		io.ReadFull(s, buf)
		offset -= length
	}
	return offset, nil
}
