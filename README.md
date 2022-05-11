# distributed-storage-system
go语言实现的简单分布式存储系统

系统适用于Linux系统，部署后可以通过接口访问，后续有时间会开发一个客户端。
因为采用的数据冗余技术问题，所以需要保证至少有六台存储设备。

### 主要功能

1. 文件上传:一次性上传所有文件和多次分段上传两种方式
2. 文件下载：可以选择下载文件所有内容和部分内容，下载结果可以选择压缩版本或者未压缩版本
3. 文件版本查询：因为使用了版本管理，会记录所有上传的文件版本，可以对之间上传的版本进行查询。
4. 文件删除
5. 文件位置查询：可以查看文件存储的物理存储设备的位置。

### 配置环境

1. 设置RabbitMQ服务器

安装rabbitmq-server
```shell
sudo apt-get install rabbitmq-server
```
下载并打开rabbitmqadmin管理工具
```shell
sudo rabbitmq-plugins enable rabbitmq_management
```
```shell
wget localhost:15672/cli/rabbitmqadmin
```
创建apiServers和dataServers两个exchange
```shell
python3 rabbitmqadmin declare exchange name=apiServers type=fanout
```
```shell
python3 rabbitmqadmin declare exchange name=dataServers type=fanout
```
添加用户(账号密码自己随意设置)
```shell
sudo rabbitmqctl add_user test test
```
给用户添加权限
```shell
sudo rabbitmqctl set_permissions -p / test ".*" ".*" ".*"
```

2. 设置elasticsearch元数据区

安装elasticsearch可以根据这个链接进行
https://github.com/1135470453/notes/blob/master/elasticsearch/elasticsearch.md

创建index

PUT IP:9200/metadata

创建type

POST IP:9200/metadata/objects
```json
{
	"mappings":{
		"objects":{
			"properties":{
                "name": {
                    "type": "text",
                    "fields": {
                      "keyword": {
                        "type": "keyword"
                      }
                    }
                  },
				"version":{
					"type":"integer"
				},
				"size":{
					"type":"integer"
				},
				"hash":{
					"type":"string"
				}
			}	
		}
	}
}
```
需要先随意添加一个数据才可以使用

例如： POST http://ip:9200/metadata/objects/test3_1?op_type=create
```json
{
  "name":"test3",
  "version":1,
  "size":13,
  "hash":"2oUvHeq7jQ27Va2y/usI1kSX4cETY9LuevZU9RT+Fuc="
}
```
### 搭建系统
系统分为两个部分，分别为处理用户请求的apiServer和负责数据存储的dataServer。
在运行以上任意一个server时，需要首先在运行的服务器设置rabbitmq和elasticsearch的环境变量
```shell
export RABBITMQ_SERVER=amqp://test:test@ip:5672/
```
```shell
export ES_SERVER=ip:9200
```
##### apiServer
apiServer负责接收用户对系统的请求，apiServer启动后，用户向该apiServer发送请求。
apiServer之间并行处理用户请求。

启动apiServer
```shell
LISTEN_ADDRESS=ip:port go run apiServer/apiServer.go
```

##### dataServer
dataServer负责数据存储业务，启动dataServer前需要首先创建存储文件的文件夹，文件夹下有负责存储临时文件的temp文件夹和负责存储正式文件的objects文件夹

启动dataServer,/storage为创建的文件夹
```shell
LISTEN_ADDRESS=ip:port STORAGE_ROOT=/storage go run dataServer/dataServer.go
```

### 使用
通过向apiServer发送请求的方法使用，假设apiServer外网地址为10.0.2.1:8082

通过该指令获取文件SHA-256编码，假设为hash
```shell
openssl dgst -sha256 -binary  | base64
```
***单次上传所有文件***

PUT方法   

10.0.2.1:8082/objects/filename(filename为自己定义的文件名称)

头部： {"key":"Digest","value":"hash"}

附带文件

成功返回200

***分批次上传文件***

发送文件信息，创建该文件在服务器的环境

POST 10.0.2.1:8082/objects/filename(filename为自己定义的文件名称)
    
头部：{"key":"Digest","value":"SHA-256="hash"},{"key":"Size","value":"文件大小"}

成功返回201和token

查询已经上传的大小

HEAD 10.0.2.1:8082/temp/token

返回信息的头部中包含已经上传的数据大小

分段上传数据

PUT 10.0.2.1:8082/temp/token

头部可选，不选则从服务器中文件的0字节处继续上传

头部：{"key":"range","value":"bytes="start-"}(start为继续上传的位置)

***下载文件***

GET 10.0.2.1:8082/objects/filename

头部可选，写则下载压缩格式，不写则下载不压缩格式

头部{"key":"Accept_Encoding","value":"gzip"}

***查看文件信息***

filename可选，不写则返回所有文件信息。version可选，不选则所有版本

GET 10.0.2.1:8082/versions/filename?version=..

***查看文件存储位置***

hash中/要改为%2F

GET 10.0.2.1:8082/locate/hash

***删除文件***

DELETE 10.0.2.1:8082/objects/filename

### 测试
因为测试时只有一台服务器，所以使用ifconfig创建别名的方式来实现
```shell
ifconfig thh0:1 10.0.1.1/16
ifconfig thh0:1 10.0.1.2/16
ifconfig thh0:1 10.0.1.3/16
ifconfig thh0:1 10.0.1.4/16
ifconfig thh0:1 10.0.1.5/16
ifconfig thh0:1 10.0.1.6/16
```
启动一个apiServer
```shell
LISTEN_ADDRESS=内网ip:8082 go run apiServer/apiServer.go  >log/apilog1 2>&1 &
```

启动六个dataServer
```shell
LISTEN_ADDRESS=10.0.1.1:8081 STORAGE_ROOT=/xiangmu/distributedStorageSystem/server/store/1 go run dataServer/dataServer.go >log/datalog1 2>&1 &
```

```shell
LISTEN_ADDRESS=10.0.1.2:8081 STORAGE_ROOT=/xiangmu/distributedStorageSystem/server/store/2 go run dataServer/dataServer.go >log/datalog2 2>&1 &
```

```shell
LISTEN_ADDRESS=10.0.1.3:8081 STORAGE_ROOT=/xiangmu/distributedStorageSystem/server/store/3 go run dataServer/dataServer.go >log/datalog3 2>&1 &
```

```shell
LISTEN_ADDRESS=10.0.1.4:8081 STORAGE_ROOT=/xiangmu/distributedStorageSystem/server/store/4 go run dataServer/dataServer.go >log/datalog4 2>&1 &
```

```shell
LISTEN_ADDRESS=10.0.1.5:8081 STORAGE_ROOT=/xiangmu/distributedStorageSystem/server/store/5 go run dataServer/dataServer.go >log/datalog5 2>&1 &
```

```shell
LISTEN_ADDRESS=10.0.1.6:8081 STORAGE_ROOT=/xiangmu/distributedStorageSystem/server/store/6 go run dataServer/dataServer.go >log/datalog6 2>&1 &
```
***单次上传所有文件***

上传名为background.png的文件，查看文件的编码为js0SsjuvtzK/YuSq3mZawJv/2RtAiRPreR+WJDcWG3o=

发送请求:

PUT 外网ip:8082/objects/background.png 

头部：

{"key":"Digest","value":"SHA-256=js0SsjuvtzK/YuSq3mZawJv/2RtAiRPreR+WJDcWG3o="}


返回200，上传成功

***分批次上传文件***
上传文件hash和大小，大小为100000字节
POST http://外网ip:8082/objects/10000test

头部
[{"key":"Digest","value":"SHA-256=/r6Oc8tjb0wRC1yeBPRUWNQnJn1IHgHhrJF6GMMn7W0=","description":""},{"key":"Size","value":"100000","description":""}]

返回201 created和/temp/eyJOYW1lIjoiZm9yVGVzdCIsIlNpemUiOjEwMDAwMCwiSGFzaCI6IkNKMERLblNIYWNqWDhwK2ptUTBzNStsUmxnU2xOSTJxaEZrd0laNk85QTg9IiwiU2VydmVycyI6WyIxMC4wLjEuNDo4MDgxIiwiMTAuMC4xLjM6ODA4MSIsIjEwLjAuMS4xOjgwODEiLCIxMC4wLjEuMjo4MDgxIiwiMTAuMC4xLjY6ODA4MSIsIjEwLjAuMS41OjgwODEiXSwiVXVpZHMiOlsiYWY2MWIxMTMtMWQ3YS00ODZjLWFkMjUtNjVhNjg3ODFhNjY0IiwiNTc2MWVlYTktZWI3Ny00YjZjLWJjZGEtZmI2ZDkzMDcwYzY5IiwiM2I1NjJhMTctMjZmNy00Mjg5LWJiYmItYjRlMzEwMzk1ZjhmIiwiNGEwMDc2YjctMDMwZi00ZTNhLTk3OTMtZGRiNDUxOWFhOWY2IiwiYWRkOTAyMDgtOGJlNC00M2EyLWFkMmMtMGYyMGMxNWQ1ZjFmIiwiMzJiNjFlNGUtNjk2OC00NzdmLTk5OWEtYTIxNzgzOGM5NDZlIl19

其中/temp/之后的字符为token

查询已经上传的数据大小

HEAD http://外网ip:8082/temp/eyJOYW1lIjoiZm9yVGVzdCIsIlNpemUiOjEwMDAwLCJIYXNoIjoiQ0owREtuU0hhY2pYOHAram1RMHM1K2xSbGdTbE5JMnFoRmt3SVo2TzlBOD0iLCJTZXJ2ZXJzIjpbIjEwLjAuMS41OjgwODEiLCIxMC4wLjEuMTo4MDgxIiwiMTAuMC4xLjQ6ODA4MSIsIjEwLjAuMS4yOjgwODEiLCIxMC4wLjEuMzo4MDgxIiwiMTAuMC4xLjY6ODA4MSJdLCJVdWlkcyI6WyI2ODBmNDU4MC1jMzYwLTQ5YzQtYWNmMi04OGRhMjkyMWExYWEiLCJhNmNlMmE3Mi02OGU2LTQxZDgtYmMzMS03ODdkNDlmYmQxNjgiLCI5ZDFmNjY2OS04MWNjLTRlMzUtYjNjZS04YTg1ZTJlNjI3MmQiLCI1ZDdjOGNhYS0zZjViLTQzYjYtOTc3Zi05MjQ5OTVjMjVlMDQiLCI0OTNiZjczMi04YjJhLTQ2NWUtYTYwNi03OGE2ZDUwNTE0ODIiLCIxYjhkM2YyMi1kOTczLTQ1ZDUtOTZlNi00NjA5MjhhZTBlNTUiXX0=

返回200,content-length:0

首先上传50000字节长度的数据

PUT http://外网ip:8082/temp/eyJOYW1lIjoiZm9yVGVzdCIsIlNpemUiOjEwMDAwMCwiSGFzaCI6IkNKMERLblNIYWNqWDhwK2ptUTBzNStsUmxnU2xOSTJxaEZrd0laNk85QTg9IiwiU2VydmVycyI6WyIxMC4wLjEuNDo4MDgxIiwiMTAuMC4xLjM6ODA4MSIsIjEwLjAuMS4xOjgwODEiLCIxMC4wLjEuMjo4MDgxIiwiMTAuMC4xLjY6ODA4MSIsIjEwLjAuMS41OjgwODEiXSwiVXVpZHMiOlsiYWY2MWIxMTMtMWQ3YS00ODZjLWFkMjUtNjVhNjg3ODFhNjY0IiwiNTc2MWVlYTktZWI3Ny00YjZjLWJjZGEtZmI2ZDkzMDcwYzY5IiwiM2I1NjJhMTctMjZmNy00Mjg5LWJiYmItYjRlMzEwMzk1ZjhmIiwiNGEwMDc2YjctMDMwZi00ZTNhLTk3OTMtZGRiNDUxOWFhOWY2IiwiYWRkOTAyMDgtOGJlNC00M2EyLWFkMmMtMGYyMGMxNWQ1ZjFmIiwiMzJiNjFlNGUtNjk2OC00NzdmLTk5OWEtYTIxNzgzOGM5NDZlIl19
附带长度50000字节长度数据
返回200

查询已经上传的数据大小，返回32000.这里应该是50000，但是只能上传32000，这里在系统设计上有些问题。

PUT http://外网ip:8082/temp/eyJOYW1lIjoiZm9yVGVzdCIsIlNpemUiOjEwMDAwMCwiSGFzaCI6IkNKMERLblNIYWNqWDhwK2ptUTBzNStsUmxnU2xOSTJxaEZrd0laNk85QTg9IiwiU2VydmVycyI6WyIxMC4wLjEuNDo4MDgxIiwiMTAuMC4xLjM6ODA4MSIsIjEwLjAuMS4xOjgwODEiLCIxMC4wLjEuMjo4MDgxIiwiMTAuMC4xLjY6ODA4MSIsIjEwLjAuMS41OjgwODEiXSwiVXVpZHMiOlsiYWY2MWIxMTMtMWQ3YS00ODZjLWFkMjUtNjVhNjg3ODFhNjY0IiwiNTc2MWVlYTktZWI3Ny00YjZjLWJjZGEtZmI2ZDkzMDcwYzY5IiwiM2I1NjJhMTctMjZmNy00Mjg5LWJiYmItYjRlMzEwMzk1ZjhmIiwiNGEwMDc2YjctMDMwZi00ZTNhLTk3OTMtZGRiNDUxOWFhOWY2IiwiYWRkOTAyMDgtOGJlNC00M2EyLWFkMmMtMGYyMGMxNWQ1ZjFmIiwiMzJiNjFlNGUtNjk2OC00NzdmLTk5OWEtYTIxNzgzOGM5NDZlIl19

头部[{"key":"range","value":"bytes=32000-","description":""}]

附带长度68000字节长度数据

返回200

***下载数据***

GET http://外网ip:8082/objects/background.png 

[{"key":"Accept-Encoding","value":"gzip","description":""}]

得到之前上传到图片

***查看文件信息***

GET http://外网ip:8082/versions/background.png

返回：{"Name":"background.png","Version":1,"Size":1875598,"Hash":"js0SsjuvtzK/YuSq3mZawJv/2RtAiRPreR+WJDcWG3o="}

***查看文件存储位置***

查看background.png位置

GET http://外网ip:8082/locate/js0SsjuvtzK%2FYuSq3mZawJv%2F2RtAiRPreR+WJDcWG3o=

返回{"0":"10.0.1.2:8081","1":"10.0.1.3:8081","2":"10.0.1.1:8081","3":"10.0.1.6:8081","4":"10.0.1.4:8081","5":"10.0.1.5:8081"}

***删除文件***

http://外网ip:8082/objects/background.png 

返回200
