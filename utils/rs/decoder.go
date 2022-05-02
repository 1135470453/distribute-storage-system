package rs

import (
	"github.com/klauspost/reedsolomon"
	"io"
	"log"
)

type decoder struct {
	readers   []io.Reader         //用于读取保存完好的分片
	writers   []io.Writer         //用于写入恢复的新分片
	enc       reedsolomon.Encoder //用于RS解码
	size      int64               //文件大小
	cache     []byte              //缓存
	cacheSize int
	total     int64 //当前已读的字节
}

func NewDecoder(readers []io.Reader, writers []io.Writer, size int64) *decoder {
	enc, _ := reedsolomon.New(DATA_SHARDS, PARITY_SHARDS)
	return &decoder{readers, writers, enc, size, nil, 0, 0}
}

//调用getData. 读取缓存中的数据
func (d *decoder) Read(p []byte) (n int, err error) {
	log.Println("decoder read start")
	//缓存中没有数据
	if d.cacheSize == 0 {
		//获取数据
		e := d.getData()
		//表示已经全部获取
		if e != nil {
			return 0, e
		}
	}
	length := len(p)
	if d.cacheSize < length {
		length = d.cacheSize
	}
	d.cacheSize -= length
	//将缓存中length长度的数据复制给p
	copy(p, d.cache[:length])
	//将缓存中已经复制的内容删除
	d.cache = d.cache[length:]
	log.Println("decoder read end")
	return length, nil
}

//读取decoder中reader和writer的数据并保留在缓存中,并使用enc修复损坏数据
func (d *decoder) getData() error {
	log.Println("decoder getData start")
	//判断已经解码的数据和文件的总大小,若相等则表示已经全部读取
	if d.total == d.size {
		return io.EOF
	}
	//保存相应分片中读取的数据
	shards := make([][]byte, ALL_SHARDS)
	repairIds := make([]int, 0)
	for i := range shards {
		if d.readers[i] == nil { //分片已丢失,保存在repairIds中
			repairIds = append(repairIds, i)
		} else { //从分片中读取至多八千字节内容
			shards[i] = make([]byte, BLOCK_PER_SHARD)
			n, e := io.ReadFull(d.readers[i], shards[i])
			if e != nil && e != io.EOF && e != io.ErrUnexpectedEOF {
				shards[i] = nil
			} else if n != BLOCK_PER_SHARD {
				shards[i] = shards[i][:n]
			}
		}
	}
	//恢复分片
	e := d.enc.Reconstruct(shards)
	if e != nil {
		return e
	}
	//将恢复的分片写入对应的server
	for i := range repairIds {
		id := repairIds[i]
		//这里调用TempPutStream.write方法写入
		d.writers[id].Write(shards[id])
	}
	//遍历数据分片,将数据添加到缓存中
	for i := 0; i < DATA_SHARDS; i++ {
		shardSize := int64(len(shards[i]))
		if d.total+shardSize > d.size {
			shardSize -= d.total + shardSize - d.size
		}
		d.cache = append(d.cache, shards[i][:shardSize]...)
		d.cacheSize += int(shardSize)
		d.total += shardSize
	}
	log.Println("decoder getData end")
	return nil
}
