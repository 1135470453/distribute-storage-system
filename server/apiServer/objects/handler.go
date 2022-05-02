package objects

import (
	"log"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	log.Println("apiServer get a handler")
	m := r.Method
	/*
		POST /objects/filename
		Digest: SHA-256=<hash>
		size:<size>
		返回对应的token，用于后续上传
	*/
	if m == http.MethodPost {
		log.Println("apiServer get a post handler")
		post(w, r)
		return
	}
	/*
		PUT /objects/fileName
		head: "Digest: SHA-256=<hash>"

		将数据直接存储到dataServer,在元数据中记录
	*/
	if m == http.MethodPut {
		log.Println("apiServer get a put handler")
		put(w, r)
		return
	}
	//GET /objects/fileName?version=..
	//range:byte=first-
	//Accept-Encoding:gzip
	//获取指定文件名的文件
	//range表示偏移量
	//Accept-Encoding表示获取gzip压缩后的对象内容，有就获取gzip压缩文件，无则获取源文件
	if m == http.MethodGet {
		log.Println("apiServer get a get handler")
		get(w, r)
		return
	}
	/*
		通过增加新版本并且size和hash设置为0的方法删除
	*/
	if m == http.MethodDelete {
		log.Println("apiServer get a del handler")
		del(w, r)
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}
