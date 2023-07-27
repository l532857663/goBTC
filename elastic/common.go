package elastic

import (
	"encoding/json"
	"goBTC/utils/http"
	"path/filepath"
)

var (
	elasticHost = "http://localhost:9200"
	defaultType = "_doc"
	username    = ""
	password    = ""
)

func InitElasticInfo(conf ElasticConfig) {
	elasticHost = conf.Host
	username = conf.Username
	password = conf.Password
}

type UrlFilter struct {
	Index  string
	Type   string
	Id     string
	Action string
}

func GetElasticUrl(filter UrlFilter) string {
	if filter.Type == "" {
		filter.Type = defaultType
	}
	url := elasticHost + "/" + filepath.Join(filter.Index, filter.Type, filter.Id, filter.Action) // + "?pretty"
	// fmt.Printf("wch------ url: %+v\n", url)
	return url
}

func AskHttpJson(method string, filter UrlFilter, reqBody, respBody interface{}) error {
	url := GetElasticUrl(filter)
	content, _ := json.Marshal(reqBody)
	// fmt.Printf("wch------ askContent: %+v\n", string(content))
	_, body, err := http.HttpByJson(method, url, username, password, content)
	if err != nil {
		return err
	}
	// fmt.Printf("wch------ code: %+v, body: %+v\n", code, string(body))
	err = json.Unmarshal(body, respBody)
	if err != nil {
		return err
	}
	return nil
}

func AskHttp(filter UrlFilter, respBody interface{}) error {
	url := GetElasticUrl(filter)
	_, body, err := http.HttpGet(url, username, password)
	if err != nil {
		return err
	}
	// fmt.Printf("wch------ code: %+v, body: %+v\n", code, string(body))
	err = json.Unmarshal(body, respBody)
	if err != nil {
		return err
	}
	return nil
}
