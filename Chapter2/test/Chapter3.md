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

##### 启动rabbitmq
```shell
sudo rabbitmq-plugins enable rabbitmq_management
```

##### es环境变量
```shell
export ES_SERVER=101.43.155.248:9200
```
##### 获取散列值
```shell
echo -n "this object will be separate to 4+2 shards" | openssl dgst -sha256 -binary |base64
```
hash: `2oUvHeq7jQ27Va2y/usI1kSX4cETY9LuevZU9RT+Fuc=`

```shell
echo -n "this is test4" | openssl dgst -sha256 -binary |base64
```

hash: `Os/0OGFkYdCb4HxMk0iubLSAJeXOe4S1Vt/6bbNIFuU=`
##### 上传文件
```shell
curl -v 10.0.2.1:8082/objects/test5 -XPUT -d"this object will be separate to 4+2 shards" -H "Digest: SHA-256=MBMxWHrPMsuOBaVYHkwScZQRyTRMQyiKp2oelpLZza8="
```

#### 查看文件位置
```shell
curl 10.0.24.11:8082/locate/K7qv1Doqv%2F2vfwjQOmbWYBmQ6UkOUaB7w%2FHW2XUC8YA=
```

##### 添加数据
POST http://101.43.155.248:9200/metadata/objects/test3_1?op_type=create
```json
{
  "name":"test3",
  "version":1,
  "size":13,
  "hash":"2oUvHeq7jQ27Va2y/usI1kSX4cETY9LuevZU9RT+Fuc="
}
```

现存问题：


如果文件在之前没有，则会报错，无法按照version排序