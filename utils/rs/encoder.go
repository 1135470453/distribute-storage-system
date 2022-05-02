package rs

import (
	"github.com/klauspost/reedsolomon"
	"io"
	"log"
)

//用于管理编码器和对应的写入writer
type encoder struct {
	writers []io.Writer
	enc     reedsolomon.Encoder
	cache   []byte
}

//创建一个包括编码器和writer的ecoder
func NewEncoder(writers []io.Writer) *encoder {
	log.Println("NewEncoder start")
	//生成一个有4个数据片,2个校验片的编码器
	enc, _ := reedsolomon.New(DATA_SHARDS, PARITY_SHARDS)
	log.Println("NewEncoder end")
	return &encoder{writers, enc, nil}
}

//将数据写入cache缓存,缓存够大就flush
func (e *encoder) Write(p []byte) (n int, err error) {
	log.Println("encoder write start")
	length := len(p)
	current := 0
	for length != 0 {
		next := BLOCK_SIZE - len(e.cache)
		if next > length {
			next = length
		}
		e.cache = append(e.cache, p[current:current+next]...)
		if len(e.cache) == BLOCK_SIZE {
			e.Flush()
		}
		current += next
		length -= next
	}
	log.Println("encoder write end")
	return len(p), nil
}

//将cache生成RS六个分片后写入临时文件中
func (e *encoder) Flush() {
	log.Println("encoder flush start")
	if len(e.cache) == 0 {
		return
	}
	//切分为四个数据片
	shards, _ := e.enc.Split(e.cache)
	//添加两个校验片
	e.enc.Encode(shards)
	//写入文件
	for i := range shards {
		e.writers[i].Write(shards[i])
	}
	e.cache = []byte{}
}
