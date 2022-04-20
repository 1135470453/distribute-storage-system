package heartbeat

import (
	"log"
	"math/rand"
)

//随机选择一个数据服务节点
func ChooseRandomDataServer() string {
	log.Println("apiServer start ChooseRandomDataServer")
	ds := GetDataServers()
	log.Println("ds's length:" + string(len(ds)))
	n := len(ds)
	if n == 0 {
		return ""
	}
	return ds[rand.Intn(n)]
}
