package types

//用于表示文件所在数据节点的地址的数据结构
//Addr:数据节点地址
//Id:文件分片号
type LocateMessage struct {
	Addr string
	Id   int
}
