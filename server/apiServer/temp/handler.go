package temp

import "net/http"

func Handler(w http.ResponseWriter, r *http.Request) {
	m := r.Method
	/*
		HEAD /temp/token
		获取已经写入的临时文件的大小(保存在头节点中)
	*/
	if m == http.MethodHead {
		head(w, r)
		return
	}
	/*
		PUT /temp/token
		head: range: bytes=<first>-<hash>
		body: 文件内容
		将body写入token对应的临时文件
	*/
	if m == http.MethodPut {
		put(w, r)
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}
