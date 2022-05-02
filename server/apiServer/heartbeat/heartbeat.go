package heartbeat

import (
	"distributed_storage_system/utils/rabbitmq"
	"log"
	"os"
	"strconv"
	"sync"
	"time"
)

//用于存储数据服务节点: [数据服务节点]发送心跳时间
var dataServers = make(map[string]time.Time)
var mutex sync.Mutex

//监听数据服务节点的心跳,并记录到dataServers中
func ListenHeartbeat() {
	log.Println("apiServer start ListenHeartbeat")
	q := rabbitmq.New(os.Getenv("RABBITMQ_SERVER"))
	defer q.Close()
	//获取的是发送心跳的数据服务的地址
	q.Bind("apiServers")
	c := q.Consume()
	go removeExpiredDataServer()
	for msg := range c {
		dataServer, e := strconv.Unquote(string(msg.Body))
		log.Println("apiServer get " + dataServer + "" +
			"at " + time.Now().String())
		if e != nil {
			panic(e)
		}
		mutex.Lock()
		dataServers[dataServer] = time.Now()
		mutex.Unlock()
	}
	log.Println("apiServer ListenHeartbeat stop")
}

//每过五秒检查并删除上一次没有发送心跳检测的数据服务节点
func removeExpiredDataServer() {
	for {
		time.Sleep(5 * time.Second)
		mutex.Lock()
		for s, t := range dataServers {
			if t.Add(10 * time.Second).Before(time.Now()) {
				delete(dataServers, s)
			}
		}
		mutex.Unlock()
	}
}

func GetDataServers() []string {
	log.Println("apiServer start GetDataServers")
	mutex.Lock()
	defer mutex.Unlock()
	ds := make([]string, 0)
	for s, _ := range dataServers {
		ds = append(ds, s)
	}
	log.Println("apiServer GetDataServers end")
	return ds
}
