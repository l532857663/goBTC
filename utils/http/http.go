package http

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

var (
	isHttps     bool
	httpClient  *http.Client
	httpsClient *http.Client
)

func InitHttps(httpsFlage bool, caCertPath string) {
	isHttps = httpsFlage
	timeout := 100 * time.Second
	transport := &http.Transport{
		// 设置为短连接请求模式
		DisableKeepAlives: true,
	}
	if isHttps {
		// 读取根证书文件
		caCert, err := ioutil.ReadFile(caCertPath)
		if err != nil {
			fmt.Println("无法读取根证书文件:", err)
			return
		}

		// 创建根证书池
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)
		transport.TLSClientConfig = &tls.Config{
			RootCAs: caCertPool,
		}
		httpsClient = &http.Client{
			Timeout:   timeout,
			Transport: transport,
		}
	}
	// 创建自定义的 http.Client
	httpClient = &http.Client{
		Timeout:   timeout,
		Transport: transport,
	}
}

func HttpByJson(method, api, user, passwd string, data []byte) (code int, body []byte, err error) {
	client := httpClient
	if isHttps {
		client = httpsClient
	}
	req, err := http.NewRequest(method, api, bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.SetBasicAuth(user, passwd)
	res, err := client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer res.Body.Close()
	bytess, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return 0, nil, err
	}
	return res.StatusCode, bytess, nil
}

func HttpGet(url, user, passwd string) (int, []byte, error) {
	client := httpClient
	if isHttps {
		client = httpsClient
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return 0, nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(user, passwd)
	resp, err := client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, err
	}
	return resp.StatusCode, body, nil
}
