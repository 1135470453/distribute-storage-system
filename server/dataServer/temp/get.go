package temp

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

/*
GET /temp/uuid
获取uuid对应的用于保存数据的临时文件的内容
*/
func get(w http.ResponseWriter, r *http.Request) {
	//从url获取uuid
	uuid := strings.Split(r.URL.EscapedPath(), "/")[2]
	//根据uuid打开保存数据的文件
	f, e := os.Open(os.Getenv("STORAGE_ROOT") + "/temp/" + uuid + ".dat")
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	defer f.Close()
	//返回文件内容
	io.Copy(w, f)
}
