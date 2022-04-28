package locate

import (
	"distributed_storage_system/utils/rabbitmq"
	"distributed_storage_system/utils/types"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

//缓存对象
var objects = make(map[string]int)
var mutex sync.Mutex

//判断name这个地址的文件是否存在
func Locate(hash string) int {
	log.Println("Locate start")
	mutex.Lock()
	id, ok := objects[hash]
	mutex.Unlock()
	if !ok {
		log.Println("Locate end")
		return -1
	}
	log.Println("Locate end")
	return id
}

//将hash加到缓存中
func Add(hash string, id int) {
	log.Println("locate Add start")
	mutex.Lock()
	objects[hash] = id
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
		id := Locate(hash)
		if id != -1 {
			q.Send(msg.ReplyTo, types.LocateMessage{Addr: os.Getenv("LISTEN_ADDRESS"), Id: id})
		}
	}
}

//将所有文件添加到缓存中
func CollectObjects() {
	files, _ := filepath.Glob(os.Getenv("STORAGE_ROOT") + "/objects/*")
	for i := range files {
		file := strings.Split(filepath.Base(files[i]), ".")
		if len(file) != 3 {
			panic(files[i])
		}
		hash := file[0]
		id, e := strconv.Atoi(file[1])
		if e != nil {
			panic(e)
		}
		objects[hash] = id
	}
}
