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

//让dataServer建立临时文件,并返回对应的uuid
func NewTempPutStream(server, object string, size int64) (*TempPutStream, error) {
	log.Println("NewTempPutStream start")
	//post方法访问数据服务的temp接口,获得uuid
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

func (w *TempPutStream) Commit(good bool) {
	method := "DELETE"
	if good {
		method = "PUT"
	}
	request, _ := http.NewRequest(method, "http://"+w.Server+"/temp/"+w.Uuid, nil)
	client := http.Client{}
	client.Do(request)
}

func NewTempGetStream(server, uuid string) (*GetStream, error) {
	return newGetSteam("http://" + server + "/temp/" + uuid)
}
