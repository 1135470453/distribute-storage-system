##### es创建index和type
1. 创建index

PUT  IP:9200/metadata

2. 创建type

POST IP:9200/metadata/objects

```json
{
	"mappings":{
		"objects":{
			"properties":{
				"name":{
					"type":"string",
					"index":"not analyzed"
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

##### 启动rabbitmq
```shell
sudo rabbitmq-plugins enable rabbitmq_management
```

##### es环境变量
```shell
export ES_SERVER=101.43.155.248:9200
```

现存问题：
es.go 101应该改为post，否则增加新数据的时候报错

Content-Type header [] is not supported

如果文件在之前没有，则会报错，无法按照version排序