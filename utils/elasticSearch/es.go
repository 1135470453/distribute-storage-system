package elasticSearch

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

//es数据对应的结构图
type Metadata struct {
	Name    string
	Version int
	Size    int64
	Hash    string
}

type hit struct {
	Source Metadata `json:"_source"`
}

type searchResult struct {
	Hits struct {
		Total int
		Hits  []hit
	}
}

//根据name和version获取对应的元数据,并保存在metadata中返回
func getMetadata(name string, versionId int) (meta Metadata, e error) {
	log.Println("getMetadata start")
	url := fmt.Sprintf("http://%s/metadata/objects/%s_%d/_source",
		os.Getenv("ES_SERVER"), name, versionId)
	log.Println("getMetadata start Get")
	log.Println("url is " + url)
	r, e := http.Get(url)
	if e != nil {
		return
	}
	if r.StatusCode != http.StatusOK {
		e = fmt.Errorf("fail to get %s_%d:%d", name, versionId, r.StatusCode)
		return
	}
	log.Println("getMetadata get success")
	result, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(result, &meta)
	log.Println("getMetadata end")
	return
}

//获取指定名字的最新版本
func SearchLatestVersion(name string) (meta Metadata, e error) {
	log.Println("SearchLatestVersion start")
	//搜索指定名字,按照版本号降序排列并且只返回第一个(及得到该名字的最新版本)
	url := fmt.Sprintf("http://%s/metadata/_search?q=name:%s&size=1&sort=version:desc",
		os.Getenv("ES_SERVER"), url.PathEscape(name))
	log.Println("SearchLatestVersion start Get")
	log.Println("url is " + url)
	r, e := http.Get(url)
	if e != nil {
		return
	}
	if r.StatusCode != http.StatusOK {
		e = fmt.Errorf("fail to search latest metadata: %d", r.StatusCode)
		return
	}

	result, _ := ioutil.ReadAll(r.Body)
	var sr searchResult
	json.Unmarshal(result, &sr)
	if len(sr.Hits.Hits) != 0 {
		meta = sr.Hits.Hits[0].Source
		log.Println("SearchLatestVersion Get success")
	}
	log.Println("SearchLatestVersion end")
	return
}

func GetMetadata(name string, version int) (Metadata, error) {
	if version == 0 {
		return SearchLatestVersion(name)
	}
	return getMetadata(name, version)
}

//向es上传新数据
func PutMetadata(name string, version int, size int64, hash string) error {
	log.Println("PutMetadata start")
	doc := fmt.Sprintf(`{"name":"%s","version":%d,"size":%d,"hash":"%s"}`,
		name, version, size, hash)
	client := http.Client{}
	//该条数据的id用name_version表示
	url := fmt.Sprintf("http://%s/metadata/objects/%s_%d?op_type=create",
		os.Getenv("ES_SERVER"), name, version)
	log.Println("PutMetadata start put")
	//这里应该改为post
	request, _ := http.NewRequest("POST", url, strings.NewReader(doc))
	request.Header.Set("Content-Type", "application/json")
	r, e := client.Do(request)
	if e != nil {
		return e
	}
	if r.StatusCode == http.StatusConflict {
		return PutMetadata(name, version+1, size, hash)
	}
	if r.StatusCode != http.StatusCreated {
		result, _ := ioutil.ReadAll(r.Body)
		return fmt.Errorf("fail to put metadata: %d %s", r.StatusCode, string(result))
	}
	log.Println("PutMetadata put success")
	log.Println("PutMetadata end")
	return nil
}

func AddVersion(name, hash string, size int64) error {
	log.Println("AddVersion start")
	version, e := SearchLatestVersion(name)
	if e != nil {
		return e
	}
	log.Println("AddVersion end")
	return PutMetadata(name, version.Version+1, size, hash)
}

//name不为空:搜这个name的所有版本, name为空,搜所有
//from size用于分页
func SearchAllVersions(name string, from, size int) ([]Metadata, error) {
	log.Println("SearchAllVersions start")
	url := fmt.Sprintf("http://%s/metadata/_search?sort=name.keyword,version&from=%d&size=%d",
		os.Getenv("ES_SERVER"), from, size)
	if name != "" {
		url += "&q=name:" + name
	}
	log.Println("SearchAllVersions start get")
	log.Println("url is " + url)
	r, e := http.Get(url)
	if e != nil {
		return nil, e
	}
	metas := make([]Metadata, 0)
	result, _ := ioutil.ReadAll(r.Body)
	var sr searchResult
	json.Unmarshal(result, &sr)
	for i := range sr.Hits.Hits {
		metas = append(metas, sr.Hits.Hits[i].Source)
	}
	return metas, nil
}
