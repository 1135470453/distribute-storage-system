package main

import (
	"distributed_storage_system/server/apiServer/heartbeat"
	"distributed_storage_system/server/apiServer/locate"
	"distributed_storage_system/server/apiServer/objects"
	"distributed_storage_system/server/apiServer/temp"
	"distributed_storage_system/server/apiServer/versions"
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
	//见handler注释
	http.HandleFunc("/objects/", objects.Handler)
	http.HandleFunc("/temp/", temp.Handler)
	//查询文件所在的数据节点,返回数据节点的地址
	//Get	ip:port/locate/hash
	//return: {"分片号":"数据节点地址",...}
	http.HandleFunc("/locate/", locate.Handler)
	//查询文件的版本信息
	//Get	ip:port/versions/fileName
	//return: {"Name":"fileName","Version":20,"Hash":"hash"}
	http.HandleFunc("/versions/", versions.Handler)
	log.Fatal(http.ListenAndServe(os.Getenv("LISTEN_ADDRESS"), nil))
}
