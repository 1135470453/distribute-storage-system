package main

import (
	"distributed_storage_system/Chapter2/dataServer/heartbeat"
	"distributed_storage_system/Chapter2/dataServer/locate"
	"distributed_storage_system/Chapter2/dataServer/objects"
	"distributed_storage_system/Chapter2/dataServer/temp"
	"log"
	"net/http"
	"os"
)

//数据服务层主函数
func main() {
	//每隔5s向apiServers发送心跳信号，信号内容为自己的地址
	log.Println("dataServer start")
	//在程序启动时候扫描磁盘,之后定位不再需要再次访问磁盘,只需要搜索内存
	locate.CollectObjects()
	go heartbeat.StartHeartbeat()
	/*连接dataservers,如果收到来自接口层的文件请求,就判断自己的存储中有无
	若有,则给这个对应的接口层返回自己的地址
	*/
	go locate.StartLocate()
	//接收来自接口服务层发送的文件(通过putStream的方法,在putStream.new方法中建立连接)
	http.HandleFunc("/objects/", objects.Handler)
	http.HandleFunc("/temp/", temp.Handler)
	log.Fatal(http.ListenAndServe(os.Getenv("LISTEN_ADDRESS"), nil))
}
