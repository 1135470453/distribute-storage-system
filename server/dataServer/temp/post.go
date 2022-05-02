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
	Uuid string //给文件分配的随机uuid
	Name string //文件名
	Size int64  //文件size
}

//请求格式: /temp/object
//object格式: fileHash.分片号
//创建存储文件基本信息的uuid,和用于缓存文件内容的uuid.dat,并返回uuid
func post(w http.ResponseWriter, r *http.Request) {
	log.Println("post start")
	//exec.Command执行shell命令
	//获取uuid值
	output, _ := exec.Command("uuidgen").Output()
	uuid := strings.TrimSuffix(string(output), "\n")
	//获取请求中的文件name和文件size
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
	//保存在结构体中
	t := tempInfo{uuid, name, size}
	//将结构体内容写入temp/uuid文件中
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
