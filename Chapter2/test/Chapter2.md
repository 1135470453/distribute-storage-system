
#####这是环境
```shell
export RABBITMQ_SERVER=amqp://zjw:zjw@101.43.155.248:5672/
```

#####开启dataServer
```shell
LISTEN_ADDRESS=10.0.24.11:8081 STORAGE_ROOT=/xiangmu/distributedStorageSystem/Chapter2/storage_root/1 go run dataServer/dataServer.go >log/datalog 2>&1 &
```

#####开启apiServer
```shell
LISTEN_ADDRESS=10.0.24.11:8082 go run apiServer/apiServer.go  >log/apilog 2>&1 &
```

#####put文件
```shell
curl -v 101.43.155.248:8082/objects/test2 -XPUT -d"this is object test2"
```

