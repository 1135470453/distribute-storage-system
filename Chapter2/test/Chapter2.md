
#####这是环境
```shell
export RABBITMQ_SERVER=amqp://zjw:zjw@101.43.155.248:5672/
```

#####开启dataServer
```shell
LISTEN_ADDRESS=10.0.24.11:8081 STORAGE_ROOT=/xiangmu/distributedStorageSystem/Chapter2/storage_root/1 go run dataServer/dataServer.go >log/datalog 2>&1 &
```

```shell
LISTEN_ADDRESS=10.0.1.1:8081 STORAGE_ROOT=/xiangmu/distributedStorageSystem/Chapter2/storage_root/1 go run dataServer/dataServer.go >log/datalog1 2>&1 &
```

```shell
LISTEN_ADDRESS=10.0.1.2:8081 STORAGE_ROOT=/xiangmu/distributedStorageSystem/Chapter2/storage_root/2 go run dataServer/dataServer.go >log/datalog2 2>&1 &
```

```shell
LISTEN_ADDRESS=10.0.1.3:8081 STORAGE_ROOT=/xiangmu/distributedStorageSystem/Chapter2/storage_root/3 go run dataServer/dataServer.go >log/datalog3 2>&1 &
```

```shell
LISTEN_ADDRESS=10.0.1.4:8081 STORAGE_ROOT=/xiangmu/distributedStorageSystem/Chapter2/storage_root/4 go run dataServer/dataServer.go >log/datalog4 2>&1 &
```

```shell
LISTEN_ADDRESS=10.0.1.5:8081 STORAGE_ROOT=/xiangmu/distributedStorageSystem/Chapter2/storage_root/5 go run dataServer/dataServer.go >log/datalog5 2>&1 &
```

```shell
LISTEN_ADDRESS=10.0.1.6:8081 STORAGE_ROOT=/xiangmu/distributedStorageSystem/Chapter2/storage_root/6 go run dataServer/dataServer.go >log/datalog6 2>&1 &
```
#####开启apiServer
```shell
LISTEN_ADDRESS=10.0.24.11:8082 go run apiServer/apiServer.go  >log/apilog 2>&1 &
```

```shell
LISTEN_ADDRESS=10.0.2.1:8082 go run apiServer/apiServer.go  >log/apilog1 2>&1 &
```

```shell
LISTEN_ADDRESS=10.0.2.2:8082 go run apiServer/apiServer.go  >log/apilog2 2>&1 &
```
#####put文件
```shell
curl -v 101.43.155.248:8082/objects/test2 -XPUT -d"this is object test2"
```

