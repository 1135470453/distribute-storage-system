package objectStream

import (
	"fmt"
	"io"
	"net/http"
)

//writer用于给某一dataServer发送数据
//c用于接收返回的错误
type PutStream struct {
	writer *io.PipeWriter
	c      chan error
}

//server:dataserver地址,object:文件名
func NewPutStream(server, object string) *PutStream {
	reader, writer := io.Pipe()
	//创建通道 make(chan 通道类型)
	c := make(chan error)
	go func() {
		/*reader为put的内容,reader读出的为writer写入的内容,writer由putstream保留
		接口服务层可以使用putstream向数据服务层发送文件内容Z
		*/
		request, _ := http.NewRequest("PUT", "http://"+server+"/objects/"+object, reader)
		//Client为结构体
		client := http.Client{}
		r, e := client.Do(request)
		if e == nil && r.StatusCode != http.StatusOK {
			e = fmt.Errorf("dataServer return http code %d", r.StatusCode)
		}
		c <- e
	}()
	return &PutStream{writer, c}
}

//向dataserver发送数据
func (w *PutStream) Write(p []byte) (n int, err error) {
	return w.writer.Write(p)
}

func (w *PutStream) Close() error {
	w.writer.Close()
	return <-w.c
}
