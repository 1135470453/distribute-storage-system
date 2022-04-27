package temp

import (
	"log"
	"net/http"
	"os"
	"strings"
)

//删除两个临时文件
func del(w http.ResponseWriter, r *http.Request) {
	log.Println("del start")
	uuid := strings.Split(r.URL.EscapedPath(), "/")[2]
	infoFile := os.Getenv("STORAGE_ROOT") + "/temp/" + uuid
	datFile := infoFile + ".dat"
	os.Remove(infoFile)
	os.Remove(datFile)
	log.Println("del end")
}
