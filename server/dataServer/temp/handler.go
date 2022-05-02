package temp

import (
	"log"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	log.Println("dataServer get a handler")
	m := r.Method
	/*
		HEAD /temp/uuid
		返回对应uuid.bat的文件大小
	*/
	if m == http.MethodHead {
		head(w, r)
		return
	}
	/*
		GET /temp/uuid
		获取uuid对应的用于保存数据的临时文件的内容
	*/
	if m == http.MethodGet {
		get(w, r)
		return
	}
	/*
		PUT /temp/uuid
		将临时文件添加到object下面,并删除临时文件
	*/
	if m == http.MethodPut {
		log.Println("dataServer get a put")
		put(w, r)
		return
	}
	//将文件内容写入uuid.dat,并比较size是否正确,不正确就删除文件并返回错误信息
	if m == http.MethodPatch {
		log.Println("dataServer get a patch")
		patch(w, r)
		return
	}
	//请求格式: /temp/object head：size=..
	//object格式: fileHash.分片号
	//创建存储文件基本信息的uuid文件,和用于缓存文件内容的uuid.dat文件,并返回uuid
	if m == http.MethodPost {
		log.Println("dataServer get a post")
		post(w, r)
		return
	}
	/*
		DELETE /temp/uuid
		删除两个临时文件
	*/
	if m == http.MethodDelete {
		log.Println("dataServer get a del")
		del(w, r)
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}
