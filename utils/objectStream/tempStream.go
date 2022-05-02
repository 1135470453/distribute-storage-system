package objectStream

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type TempPutStream struct {
	Server string
	Uuid   string
}

//让dataServer创建存储文件基本信息的uuid文件,和用于缓存文件内容的uuid.dat文件,并返回TempPutStream格式的uuid和server
//object格式: hash.分片号
func NewTempPutStream(server, object string, size int64) (*TempPutStream, error) {
	log.Println("NewTempPutStream start")
	//创建存储文件基本信息的uuid文件,和用于缓存文件内容的uuid.dat文件,并返回uuid
	request, e := http.NewRequest("POST", "http://"+server+"/temp/"+object, nil)
	if e != nil {
		return nil, e
	}
	request.Header.Set("size", fmt.Sprintf("%d", size))
	client := http.Client{}
	response, e := client.Do(request)
	if e != nil {
		return nil, e
	}
	//获取post返回的uuid
	uuid, e := ioutil.ReadAll(response.Body)
	if e != nil {
		return nil, e
	}
	log.Println("NewTempPutStream end")
	return &TempPutStream{server, string(uuid)}, nil
}

//根据uuid,向指定server发送patch请求,将文件上传
func (w *TempPutStream) Write(p []byte) (n int, err error) {
	log.Println("TempPutStream write start")
	request, e := http.NewRequest("PATCH", "http://"+w.Server+"/temp/"+w.Uuid, strings.NewReader(string(p)))
	if e != nil {
		return 0, e
	}
	client := http.Client{}
	r, e := client.Do(request)
	if e != nil {
		return 0, e
	}
	if r.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("dataServer return http code %d", r.StatusCode)
	}
	log.Println("TempPutStream write end")
	return len(p), nil
}

//true则将临时文件变为正式文件, false则删除临时文件
func (w *TempPutStream) Commit(good bool) {
	method := "DELETE"
	if good {
		method = "PUT"
	}
	request, _ := http.NewRequest(method, "http://"+w.Server+"/temp/"+w.Uuid, nil)
	client := http.Client{}
	client.Do(request)
}

//从dataServer获取uuid对应的用于保存文件内容的临时文件内容
func NewTempGetStream(server, uuid string) (*GetStream, error) {
	return newGetSteam("http://" + server + "/temp/" + uuid)
}
