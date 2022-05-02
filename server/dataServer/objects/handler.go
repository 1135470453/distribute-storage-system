package objects

import (
	"log"
	"net/http"
)

/*
收到请求格式:
GET /objects/filename   filename为hash.i格式
获取已经进行内容检验、被压缩的文件
*/
func Handler(w http.ResponseWriter, r *http.Request) {
	log.Println("dataServer get a handler")
	m := r.Method
	if m == http.MethodGet {
		log.Println("dataServer get a get")
		get(w, r)
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}
