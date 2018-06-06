package utils

import (
	"io/ioutil"
	"net/http"
	"strings"
)

type Request struct {
	Data    string            `json:"data"`
	Headers map[string]string `json:"headers"`
	Url     string            `json:"url"`
	Method  string            `json:"method"`
}

type H map[string]interface{}

func (r *Request) Http() H {
	req, err := http.NewRequest(r.Method, r.Url, strings.NewReader(r.Data))
	if err != nil {
		return H{
			"status": false,
			"result": "请求出错! Error: " + err.Error(),
		}
	}
	for k, v := range r.Headers {
		req.Header.Add(k, v)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return H{
			"status": false,
			"result": "请求出错! Error: " + err.Error(),
		}
	}
	defer resp.Body.Close()
	res, _ := ioutil.ReadAll(resp.Body)
	return H{
		"status": true,
		"result": res,
	}

}
