package locate

import (
	"distributed_storage_system/utils/rabbitmq"
	"distributed_storage_system/utils/rs"
	"distributed_storage_system/utils/types"
	"encoding/json"
	"log"
	"os"
	"time"
)

// Locate 查询文件所在的数据节点,返回数据节点的地址
func Locate(name string) (locateInfo map[int]string) {
	log.Println("apiServer locate start")
	//向dataServers exchange广播hash,dataServer收到hash,如果自己存有该文件的分片,返回自己的地址和分片号
	q := rabbitmq.New(os.Getenv("RABBITMQ_SERVER"))
	q.Publish("dataServers", name)
	c := q.Consume()
	go func() {
		time.Sleep(time.Second)
		q.Close()
	}()
	//用于存储所有的分片,[分片id]:节点地址
	locateInfo = make(map[int]string)
	for i := 0; i < rs.ALL_SHARDS; i++ {
		msg := <-c
		if len(msg.Body) == 0 {
			log.Println("apiServer locate end")
			return
		}
		var info types.LocateMessage
		json.Unmarshal(msg.Body, &info)
		locateInfo[info.Id] = info.Addr
	}
	log.Println("apiServer locate end")
	return
}

//判断是否有六个及以上节点存有分片
func Exist(name string) bool {
	log.Println("apiServer exist start")
	return len(Locate(name)) >= rs.DATA_SHARDS
}
