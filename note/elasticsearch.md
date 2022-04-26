# Chapter3中elasticsearch的设置方式
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

3. 需要先添加一个数据才可以使用

POST http://101.43.155.248:9200/metadata/objects/test3_1?op_type=create
```json
{
  "name":"test3",
  "version":1,
  "size":13,
  "hash":"2oUvHeq7jQ27Va2y/usI1kSX4cETY9LuevZU9RT+Fuc="
}
```
