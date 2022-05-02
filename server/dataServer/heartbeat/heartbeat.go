package heartbeat

import (
	"distributed_storage_system/utils/rabbitmq"
	"log"
	"os"
	"time"
)

//每隔5s向apiServers发送心跳信号，信号内容为自己的地址
func StartHeartbeat() {
	log.Println("dataServer start StartHeartbeat")
	q := rabbitmq.New(os.Getenv("RABBITMQ_SERVER"))
	defer q.Close()
	for {
		log.Println("dataServer send" + os.Getenv("LISTEN_ADDRESS") + "" +
			"at " + time.Now().String())
		q.Publish("apiServers", os.Getenv("LISTEN_ADDRESS"))
		time.Sleep(5 * time.Second)
	}
}
