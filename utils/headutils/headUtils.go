package headutils

import (
	"crypto/sha256"
	"encoding/base64"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
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

// CalculateHash 计算hash值
func CalculateHash(r io.Reader) string {
	log.Println("CalculateHash start")
	h := sha256.New()
	io.Copy(h, r)
	//h.Sum获取散列值
	//base64.StdEncoding.EncodeToString进行base64编码
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func GetOffsetFromHeader(h http.Header) int64 {
	byteRange := h.Get("range")
	if len(byteRange) < 7 {
		return 0
	}
	if byteRange[:6] != "bytes=" {
		return 0
	}
	bytePos := strings.Split(byteRange[6:], "-")
	offset, _ := strconv.ParseInt(bytePos[0], 0, 64)
	return offset
}
