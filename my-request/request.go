package my_request

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"io"
	"net"
	"net/http"
	"strings"
	"time"
)

const (
	UA      = "Mozilla/5.0 (iPad; CPU OS 11_0 like Mac OS X) AppleWebKit/604.1.34 (KHTML, like Gecko) Version/11.0 Mobile/15A5341f Safari/604.1"
	Referer = "https://www.google.com/"
)

type MyClient struct {
	http.Client
}

func NewMyClient(followRedirect bool) *MyClient {
	c := http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   30 * time.Second,
			ResponseHeaderTimeout: 30 * time.Second,
			ExpectContinueTimeout: 10 * time.Second,
		},
		Timeout: 30 * time.Second,
	}
	// https://stackoverflow.com/questions/23297520/how-can-i-make-the-go-http-client-not-follow-redirects-automatically
	// more elegant https://colobu.com/2017/04/19/go-http-redirect/
	if !followRedirect {
		c.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}

	return &MyClient{c}
}

func (c *MyClient) Request(urlStr, method, cookieStr string, body io.Reader,
	headers map[string]string, params map[string]string) (*http.Response, error) {

	method = strings.ToUpper(method)

	req, err := http.NewRequest(method, urlStr, body)
	if nil != err {
		return nil, err
	}

	hasUa, hasReferer := false, false
	for key := range headers {
		req.Header.Add(key, headers[key])
		switch strings.ToLower(key) {
		case "user-agent":
			hasUa = true
		case "referer":
			hasReferer = true
		}
	}
	if !hasUa {
		req.Header.Add("user-agent", UA)
	}
	if !hasReferer {
		req.Header.Add("referer", Referer)
	}

	if nil != params {
		q := req.URL.Query()
		for key := range params {
			q.Add(key, params[key])
		}
	}

	if "" != cookieStr {
		split := strings.Split(cookieStr, ";")
		for _, s := range split {
			cookieArr := strings.Split(strings.TrimSpace(s), "=")
			req.AddCookie(&http.Cookie{
				Name:  strings.TrimSpace(cookieArr[0]),
				Value: strings.TrimSpace(cookieArr[1]),
			})
		}
	}

	resp, err := c.Do(req)
	if nil != err {
		if err, ok := err.(net.Error); ok && err.Timeout() {
			res, err := c.retry(req)
			if nil != err {
				return nil, err
			} else {
				resp = res
			}
		} else {
			return nil, err
		}
	}

	switch resp.StatusCode {
	case http.StatusOK:
	case http.StatusFound:
		if loc := resp.Header.Get("Location"); "" != loc &&
			strings.Contains(loc, "login") {
			return resp, errors.New("Cookie maybe is invalid")
		}
	default:
		logrus.Infof("StatusCode: %d, url: %s", resp.StatusCode, urlStr)
	}

	return resp, nil
}

var (
	retry    = 0
	retryMax = 5
	interval = time.Second
)

func (c *MyClient) retry(req *http.Request) (*http.Response, error) {
	retry++
	if retry > retryMax {
		return nil, errors.New("Reach max retry times")
	}
	time.Sleep(interval)
	resp, err := c.Do(req)
	if nil != err {
		if err, ok := err.(net.Error); ok && err.Timeout() {
			return c.retry(req)
		} else {
			return resp, err
		}
	}
	return resp, nil
}
