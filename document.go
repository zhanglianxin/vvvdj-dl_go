package main

import (
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

type Documentable interface {
	getContent() *interface{}
}

type HtmlDocument struct {
	url     string
	method  string
	headers map[string]string
	params  map[string]string
}

func (hd *HtmlDocument) getContent() []byte {
	resp := makeRequest(hd.url, hd.method, hd.headers, hd.params)
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if nil != err {
		panic(err)
	}
	return b
}

type StreamDocument struct {
	url     string
	method  string
	headers map[string]string
	params  map[string]string
}

func (sd *StreamDocument) getContent() io.ReadCloser {
	resp := makeRequest(sd.url, sd.method, sd.headers, sd.params)
	// defer resp.Body.Close()
	logrus.Infof("url: [%v], length: [%v]", sd.url, resp.ContentLength)
	return resp.Body
}

func makeRequest(url string, method string, headers map[string]string, params map[string]string) *http.Response {
	method = strings.ToUpper(method)
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if nil != err {
		panic(err)
	}
	for key := range headers {
		req.Header.Add(key, headers[key])
	}
	q := req.URL.Query()
	for key := range params {
		q.Add(key, params[key])
	}
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if nil != err {
		panic(err)
	}
	return resp
}
