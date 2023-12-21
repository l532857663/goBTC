package http

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"errors"
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
		DisableKeepAlives: false,
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

func HttpGetWithHeader(url string, header map[string]string) (code int, body string, err error) {
	// 构造请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, "", errors.New("New req failed!")
	}
	// 添加header
	if len(header) > 0 {
		for k, v := range header {
			req.Header.Add(k, v)
		}
	}

	// 发送http请求
	response, err := httpClient.Do(req)
	if err != nil {
		return 0, "", err
	}
	defer response.Body.Close()

	// 读取返回body内容
	bytess, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return 0, "", err
	}
	return response.StatusCode, string(bytess), nil
}
