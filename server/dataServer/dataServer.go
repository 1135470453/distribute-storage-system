package main

import (
	"distributed_storage_system/server/dataServer/heartbeat"
	"distributed_storage_system/server/dataServer/locate"
	"distributed_storage_system/server/dataServer/objects"
	"distributed_storage_system/server/dataServer/temp"
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
	/*
		收到请求格式:
		GET /objects/filename   filename为hash.i格式
		获取已经进行内容检验、被压缩的文件
	*/
	http.HandleFunc("/objects/", objects.Handler)
	/*
		POST /temp/object head：size=..
		object格式: fileHash.分片号
		功能:创建存储文件基本信息的uuid文件,和用于缓存文件内容的uuid.dat文件,并返回uuid
	*/
	http.HandleFunc("/temp/", temp.Handler)
	log.Fatal(http.ListenAndServe(os.Getenv("LISTEN_ADDRESS"), nil))
}
