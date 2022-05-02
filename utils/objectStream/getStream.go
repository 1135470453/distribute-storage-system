package objectStream

import (
	"fmt"
	"io"
	"net/http"
)

type GetStream struct {
	reader io.Reader
}

//向dataServer发送get请求,返回所获得的数据
func newGetSteam(url string) (*GetStream, error) {
	//获取已经进行内容检验、被压缩的文件
	r, e := http.Get(url)
	if e != nil {
		return nil, e
	}
	if r.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("dataServer return http code %d", r.StatusCode)
	}
	//以GetStream形式返回file
	return &GetStream{r.Body}, nil
}

/*
server:存有object的dataServer地址
object:hash.i格式,为存在dataServer的文件名称
//向Server对应的dataServe发出请求，获取以GetStream形式的已经进行内容检验、被压缩的object
*/
func NewGetSteam(server, object string) (*GetStream, error) {
	if server == "" || object == "" {
		return nil, fmt.Errorf("invalid server %s object %s", server, object)
	}
	return newGetSteam("http://" + server + "/objects/" + object)
}

func (r *GetStream) Read(p []byte) (n int, err error) {
	return r.reader.Read(p)
}
