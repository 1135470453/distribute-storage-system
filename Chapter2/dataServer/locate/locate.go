package locate

import (
	"distributed_storage_system/utils/rabbitmq"
	"os"
	"strconv"
)

//判断name这个地址的文件是否存在
func Locate(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

/*连接dataservers,如果收到来自接口层的文件请求,就判断自己的存储中有无
若有,则给这个对应的接口层返回自己的地址
*/
func StartLocate() {
	q := rabbitmq.New(os.Getenv("RABBITMQ_SERVER"))
	defer q.Close()
	q.Bind("dataServers")
	c := q.Consume()
	for msg := range c {
		object, e := strconv.Unquote(string(msg.Body))
		if e != nil {
			panic(e)
		}
		if Locate(os.Getenv("STORAGE_ROOT") + "/objects/" + object) {
			q.Send(msg.ReplyTo, os.Getenv("LISTEN_ADDRESS"))
		}
	}
}
