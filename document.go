package main

import (
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
	body := makeRequest(hd.url, hd.method, hd.headers, hd.params)
	defer (*body).Close()
	b, err := ioutil.ReadAll(*body)
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

func (sd *StreamDocument) getContent() *io.ReadCloser {
	body := makeRequest(sd.url, sd.method, sd.headers, sd.params)
	defer (*body).Close()
	return body
}

func makeRequest(url string, method string, headers map[string]string, params map[string]string) *io.ReadCloser {
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

	if http.StatusOK != resp.StatusCode {
		return nil
	} else {
		return &resp.Body
	}
}
