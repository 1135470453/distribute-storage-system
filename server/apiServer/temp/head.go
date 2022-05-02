package temp

import (
	"distributed_storage_system/utils/rs"
	"fmt"
	"log"
	"net/http"
	"strings"
)

/*
HEAD /temp/token
获取已经写入的临时文件的大小(保存在头节点中)
*/
func head(w http.ResponseWriter, r *http.Request) {
	//从url获取token
	token := strings.Split(r.URL.EscapedPath(), "/")[2]
	//根据token生成RSResumablePutStream，用于写入数据
	stream, e := rs.NewRSResumablePutStreamFromToken(token)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusForbidden)
		return
	}
	//获取当前已经写入的大小
	current := stream.CurrentSize()
	if current == -1 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	//将已经写入的大小放入头节点返回
	w.Header().Set("content-length", fmt.Sprintf("%d", current))
}
