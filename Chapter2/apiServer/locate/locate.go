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

func Locate(name string) (locateInfo map[int]string) {
	log.Println("apiServer locate start")
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
			return
		}
		var info types.LocateMessage
		json.Unmarshal(msg.Body, &info)
		locateInfo[info.Id] = info.Addr
	}
	log.Println("apiServer locate end")
	return
}

func Exist(name string) bool {
	log.Println("apiServer exist start")
	return len(Locate(name)) >= rs.DATA_SHARDS
}
