package locate

import (
	"distributed_storage_system/utils/rabbitmq"
	"os"
	"strconv"
	"time"
)

// Locate 借助dataServers exchange来获取哪一个数据节点含有这个数据,并返回这个数据节点的地址
func Locate(name string) string {
	q := rabbitmq.New(os.Getenv("RABBITMQ_SERVER"))
	q.Publish("dataServers", name)
	c := q.Consume()
	go func() {
		time.Sleep(time.Second)
		q.Close()
	}()
	msg := <-c
	s, _ := strconv.Unquote(string(msg.Body))
	return s
}

// Exist 判断是否存在这个文件,但是没有使用到
func Exist(name string) bool {
	return Locate(name) != ""
}
