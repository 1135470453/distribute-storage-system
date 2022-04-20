package main

import (
	"distributed_storage_system/Chapter2/apiServer/heartbeat"
	"distributed_storage_system/Chapter2/apiServer/locate"
	"distributed_storage_system/Chapter2/apiServer/objects"
	"log"
	"net/http"
	"os"
)

func main() {
	/*监听数据服务节点的心跳,并记录到dataServers中,
	每五秒删除上次没有发送心跳的数据服务节点
	*/
	log.Println("apiserver start")
	go heartbeat.ListenHeartbeat()
	http.HandleFunc("/objects/", objects.Handler)
	//查询文件所在的数据节点,返回数据节点的地址
	http.HandleFunc("/locate/", locate.Handler)
	log.Fatal(http.ListenAndServe(os.Getenv("LISTEN_ADDRESS"), nil))
}
