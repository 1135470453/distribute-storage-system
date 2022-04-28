package heartbeat

import (
	"log"
	"math/rand"
)

//随机选择数据服务节点
//n:需要的节点数量
//exclde:不能包含的节点(数据复原时使用,使包含正常分片的节点不再被选择)
func ChooseRandomDataServers(n int, exclude map[int]string) (ds []string) {
	log.Println("ChooseRandomDataServers start")
	candidates := make([]string, 0)
	//对exclude的键值转换,从而方便遍历
	reverseExcludeMap := make(map[string]int)
	for id, addr := range exclude {
		reverseExcludeMap[addr] = id
	}
	servers := GetDataServers()
	for i := range servers {
		s := servers[i]
		_, excluded := reverseExcludeMap[s]
		if !excluded {
			candidates = append(candidates, s)
		}
	}
	length := len(candidates)
	if length < n {
		return
	}
	//随机选择节点
	p := rand.Perm(length)
	for i := 0; i < n; i++ {
		ds = append(ds, candidates[p[i]])
	}
	log.Println("ChooseRandomDataServers end")
	return
}
