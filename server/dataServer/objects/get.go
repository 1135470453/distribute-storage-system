package objects

import (
	"compress/gzip"
	"crypto/sha256"
	"distributed_storage_system/server/dataServer/locate"
	"encoding/base64"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

/*
收到请求格式:
GET /objects/filename   filename为hash.i格式
获取已经进行内容检验、被压缩的文件
*/
//将已经检验好的，并已经压缩的文件返回
func get(w http.ResponseWriter, r *http.Request) {
	log.Println("object get start")
	//获取url中的filename
	file := getFile(strings.Split(r.URL.EscapedPath(), "/")[2])
	if file == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	sendFile(w, file)
	log.Println("object get end")
}

//获取name对应的文件具体地址
//并判断文件是否损坏，若损坏则删除文件返回空值
func getFile(name string) string {
	log.Println("getFile start")
	//获取file的具体所在位置
	//file名为 整个file的hash.分片号.该分片的hash 格式存储
	files, _ := filepath.Glob(os.Getenv("STORAGE_ROOT") + "/objects/" + name + ".*")
	if len(files) != 1 {
		return ""
	}
	file := files[0]
	//h可以计算出读入的文件的hash值
	h := sha256.New()
	//将file读入h
	sendFile(h, file)
	//计算文件的hash值
	d := url.PathEscape(base64.StdEncoding.EncodeToString(h.Sum(nil)))
	//获取文件名中该分片的hash
	hash := strings.Split(file, ".")[2]
	//比较hash值来判断文件是否损坏
	//若损坏则删除文件
	if d != hash {
		log.Println("object hash mismatch, remove", file)
		locate.Del(hash)
		os.Remove(file)
		return ""
	}
	log.Println("getFile end")
	return file
}

//将被gzip压缩的文件file传入到w中
func sendFile(w io.Writer, file string) {
	//打开文件
	f, e := os.Open(file)
	if e != nil {
		log.Println(e)
		return
	}
	defer f.Close()
	//gzipStream为指向gzip.Reader的指针

	gzipStream, e := gzip.NewReader(f)
	if e != nil {
		log.Println(e)
		return
	}
	//从gizipStream获取文件时,文件会被先被gzip压缩，然后被取出来
	io.Copy(w, gzipStream)
	gzipStream.Close()
}
