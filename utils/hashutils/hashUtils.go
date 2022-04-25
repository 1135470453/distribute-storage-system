package hashutils

import (
	"log"
	"net/http"
	"strconv"
)

//获取hash值
func GetHashFromHeader(h http.Header) string {
	log.Println("GetHashFromHeader start")
	digest := h.Get("digest")
	log.Println("digest is " + digest)
	if len(digest) < 9 {
		return ""
	}
	if digest[:8] != "SHA-256=" {
		return ""
	}
	log.Println("GetHashFromHeader end")
	return digest[8:]
}

func GetSizeFromHeader(h http.Header) int64 {
	log.Println("GetSizeFromHeader start")
	size, _ := strconv.ParseInt(h.Get("content-length"), 0, 64)
	log.Printf("size is %d", size)
	log.Println("GetSizeFromHeader end")
	return size
}
