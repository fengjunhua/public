package httputils

import (
	"bytes"
	"encoding/json"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"time"
)

/*
1.通用请求操作,适用于GET,POST,PUT,DELETE,PATCH,HEAD等操作,GET,DELETE,HEAD等操作参数传递的时候需要将params和data
  设置为nil
*/
type Client http.Client

func (client *Client) Request(method string, url string, params map[string]string, headers map[string]string, data []byte) {

}

func HttpDo(method, url string, headers, params map[string]string, data []byte) (interface{}, bool) {

	//创建request对象
	req, err := http.NewRequest(method, url, bytes.NewBuffer(data))
	if err != nil {
		return errors.New("new request is fail: %v \n"), false
	}

	// add headers
	req.Header.Set("Content-type", "application/json")
	if headers != nil {
		for key, value := range headers {
			req.Header.Add(key, value)
		}
	}

	//add params
	query := req.URL.Query()
	if params != nil {
		for key, value := range params {
			query.Add(key, value)
		}
		req.URL.RawQuery = query.Encode()
	}

	//默认的resp ,err :=  http.DefaultClient.Do(req)
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err2 := client.Do(req)
	if err2 != nil {
		return errors.New("new request is fail: %v \n"), false
	}
	defer resp.Body.Close()

	//解析返回的数据
	body, err3 := ioutil.ReadAll(resp.Body)

	temp := make(map[string]interface{}, 0)

	err = json.Unmarshal(body, &temp)

	if err3 != nil {

		return errors.New("new request is fail: %v \n"), false
	}

	return temp, true

}

/*
 */

func HttpGet(method, url string, headers, params map[string]string) (interface{}, bool) {

	//创建request请求对象
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return errors.New("new request is fail: %v \n"), false
	}

	// add headers
	req.Header.Set("Content-type", "application/json")
	if headers != nil {
		for key, value := range headers {
			req.Header.Add(key, value)
		}
	}

	//add params
	query := req.URL.Query()
	if params != nil {
		for key, value := range params {
			query.Add(key, value)
		}
		req.URL.RawQuery = query.Encode()
	}

	//默认的resp ,err :=  http.DefaultClient.Do(req)
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err2 := client.Do(req)
	if err2 != nil {
		return errors.New("new request is fail: %v \n"), false
	}
	defer resp.Body.Close()

	//解析返回的数据
	body, err3 := ioutil.ReadAll(resp.Body)

	temp := make(map[string]interface{}, 0)

	err = json.Unmarshal(body, &temp)

	if err3 != nil {

		return errors.New("new request is fail: %v \n"), false
	}

	return temp, true

}

func HttpPost(method, url string, headers, params map[string]string, data []byte) (interface{}, bool) {

	//创建request请求对象
	req, err := http.NewRequest(method, url, bytes.NewBuffer(data))
	if err != nil {
		return errors.New("new request is fail: %v \n"), false
	}

	// add headers
	req.Header.Set("Content-type", "application/json")
	if headers != nil {
		for key, value := range headers {
			req.Header.Add(key, value)
		}
	}

	//add params
	query := req.URL.Query()
	if params != nil {
		for key, value := range params {
			query.Add(key, value)
		}
		req.URL.RawQuery = query.Encode()
	}

	//默认的resp ,err :=  http.DefaultClient.Do(req)
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err2 := client.Do(req)
	if err2 != nil {
		return errors.New("new request is fail: %v \n"), false
	}
	defer resp.Body.Close()

	//解析返回的数据
	body, err3 := ioutil.ReadAll(resp.Body)

	temp := make(map[string]interface{}, 0)

	err = json.Unmarshal(body, &temp)

	if err3 != nil {

		return errors.New("new request is fail: %v \n"), false
	}

	return temp, true
}

func HttpPut(method, url string, headers, params map[string]string, data []byte) (interface{}, bool) {

	//创建request请求对象
	req, err := http.NewRequest(method, url, bytes.NewBuffer(data))
	if err != nil {
		return errors.New("new request is fail: %v \n"), false
	}

	// add headers
	req.Header.Set("Content-type", "application/json")
	if headers != nil {
		for key, value := range headers {
			req.Header.Add(key, value)
		}
	}

	//add params
	query := req.URL.Query()
	if params != nil {
		for key, value := range params {
			query.Add(key, value)
		}
		req.URL.RawQuery = query.Encode()
	}

	//默认的resp ,err :=  http.DefaultClient.Do(req)
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err2 := client.Do(req)
	if err2 != nil {
		return errors.New("new request is fail: %v \n"), false
	}
	defer resp.Body.Close()

	//解析返回的数据
	body, err3 := ioutil.ReadAll(resp.Body)

	temp := make(map[string]interface{}, 0)

	err = json.Unmarshal(body, &temp)

	if err3 != nil {

		return errors.New("new request is fail: %v \n"), false
	}

	return temp, true

}

func HttpDelete(method, url string, headers map[string]string) (interface{}, bool) {

	//创建request请求对象
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return errors.New("new request is fail: %v \n"), false
	}

	// add headers
	req.Header.Set("Content-type", "application/json")
	if headers != nil {
		for key, value := range headers {
			req.Header.Add(key, value)
		}
	}

	//默认的resp ,err :=  http.DefaultClient.Do(req)
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err2 := client.Do(req)
	if err2 != nil {
		return errors.New("new request is fail: %v \n"), false
	}
	defer resp.Body.Close()

	return nil, true
}

func HttpPatch(method, url string, headers, params map[string]string, data []byte) (interface{}, bool) {
	//创建request请求对象
	req, err := http.NewRequest(method, url, bytes.NewBuffer(data))
	if err != nil {
		return errors.New("new request is fail: %v \n"), false
	}

	// add headers
	req.Header.Set("Content-type", "application/json")
	if headers != nil {
		for key, value := range headers {
			req.Header.Add(key, value)
		}
	}

	//add params
	query := req.URL.Query()
	if params != nil {
		for key, value := range params {
			query.Add(key, value)
		}
		req.URL.RawQuery = query.Encode()
	}

	//默认的resp ,err :=  http.DefaultClient.Do(req)
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err2 := client.Do(req)
	if err2 != nil {
		return errors.New("new request is fail: %v \n"), false
	}
	defer resp.Body.Close()

	//解析返回的数据
	body, err3 := ioutil.ReadAll(resp.Body)

	temp := make(map[string]interface{}, 0)

	err = json.Unmarshal(body, &temp)

	if err3 != nil {

		return errors.New("new request is fail: %v \n"), false
	}

	return temp, true
}

func HttpHead(method, url string, headers, params map[string]string) (interface{}, bool) {

	//创建request请求对象
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return errors.New("new request is fail: %v \n"), false
	}

	// add headers
	req.Header.Set("Content-type", "application/json")
	if headers != nil {
		for key, value := range headers {
			req.Header.Add(key, value)
		}
	}

	//add params
	query := req.URL.Query()
	if params != nil {
		for key, value := range params {
			query.Add(key, value)
		}
		req.URL.RawQuery = query.Encode()
	}

	//默认的resp ,err :=  http.DefaultClient.Do(req)
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err2 := client.Do(req)
	if err2 != nil {
		return errors.New("new request is fail: %v \n"), false
	}
	defer resp.Body.Close()

	//解析返回的数据
	body, err3 := ioutil.ReadAll(resp.Body)

	temp := make(map[string]interface{}, 0)

	err = json.Unmarshal(body, &temp)

	if err3 != nil {

		return errors.New("new request is fail: %v \n"), false
	}

	return temp, true

}
