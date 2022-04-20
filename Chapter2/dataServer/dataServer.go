package main

import (
	"distributed_storage_system/Chapter2/dataServer/heartbeat"
	"distributed_storage_system/Chapter2/dataServer/locate"
	"distributed_storage_system/Chapter2/dataServer/objects"
	"log"
	"net/http"
	"os"
)

//数据服务层主函数
func main() {
	//每隔5s向apiServers发送心跳信号，信号内容为自己的地址
	log.Println("dataServer start")
	go heartbeat.StartHeartbeat()
	/*连接dataservers,如果收到来自接口层的文件请求,就判断自己的存储中有无
	若有,则给这个对应的接口层返回自己的地址
	*/
	go locate.StartLocate()
	//接收来自接口服务层发送的文件(通过putStream的方法,在putStream.new方法中建立连接)
	http.HandleFunc("/objects/", objects.Handler)
	log.Fatal(http.ListenAndServe(os.Getenv("LISTEN_ADDRESS"), nil))
}
