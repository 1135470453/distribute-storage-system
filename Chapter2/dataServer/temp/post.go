package temp

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type tempInfo struct {
	Uuid string
	Name string
	Size int64
}

//创建存储文件基本信息的uuid,和用于缓存文件内容的uuid.dat,并返回uuid
func post(w http.ResponseWriter, r *http.Request) {
	log.Println("post start")
	//exec.Command执行shell命令
	//获取uuid值
	output, _ := exec.Command("uuidgen").Output()
	uuid := strings.TrimSuffix(string(output), "\n")
	name := strings.Split(r.URL.EscapedPath(), "/")[2]
	size, e := strconv.ParseInt(r.Header.Get("size"), 0, 64)
	log.Println("uuid :" + uuid)
	log.Println("name :" + name)
	log.Printf("size : %d", size)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	t := tempInfo{uuid, name, size}
	e = t.writeToFile()
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//创建一个uuid.bat用于存储对象的内容
	os.Create(os.Getenv("STORAGE_ROOT") + "/temp/" + t.Uuid + ".dat")
	//通过http相应返回uuid
	w.Write([]byte(uuid))
}

//将结构体内容写入temp/uuid文件中
func (t *tempInfo) writeToFile() error {
	log.Println("writeToFile start")
	f, e := os.Create(os.Getenv("STORAGE_ROOT") + "/temp/" + t.Uuid)
	if e != nil {
		return e
	}
	defer f.Close()
	b, _ := json.Marshal(t)
	f.Write(b)
	log.Println("writeToFile end")
	return nil
}
