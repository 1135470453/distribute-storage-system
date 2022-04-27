package locate

import (
	"distributed_storage_system/utils/rabbitmq"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

//缓存对象
var objects = make(map[string]int)
var mutex sync.Mutex

//判断name这个地址的文件是否存在
func Locate(hash string) bool {
	log.Println("Locate start")
	mutex.Lock()
	_, ok := objects[hash]
	mutex.Unlock()
	log.Println("Locate end")
	return ok
}

//将hash加到缓存中
func Add(hash string) {
	log.Println("locate Add start")
	mutex.Lock()
	objects[hash] = 1
	mutex.Unlock()
	log.Println("locate Add end")
}

//删除缓存中的hash
func Del(hash string) {
	log.Println("locate Del start")
	mutex.Lock()
	delete(objects, hash)
	mutex.Unlock()
	log.Println("locate Del end")
}

/*连接dataservers,如果收到来自接口层的文件请求,就判断自己的缓存中有无
若有,则给这个对应的接口层返回自己的地址
*/
func StartLocate() {
	log.Println("StartLocate start")
	q := rabbitmq.New(os.Getenv("RABBITMQ_SERVER"))
	defer q.Close()
	q.Bind("dataServers")
	c := q.Consume()
	for msg := range c {
		hash, e := strconv.Unquote(string(msg.Body))
		log.Println("hash is " + hash)
		if e != nil {
			panic(e)
		}
		exist := Locate(hash)
		if exist {
			log.Println("this server has this hash")
			q.Send(msg.ReplyTo, os.Getenv("LISTEN_ADDRESS"))
		}
	}
}

//将所有文件添加到缓存中
func CollectObjects() {
	//读取目录里的所有文件
	files, _ := filepath.Glob(os.Getenv("STORAGE_ROOT") + "/objects/*")
	for i := range files {
		hash := filepath.Base(files[i])
		objects[hash] = 1
	}
}
