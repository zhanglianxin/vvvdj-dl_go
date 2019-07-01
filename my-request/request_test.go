package my_request

import (
	"io/ioutil"
	"testing"
)

func TestNewMyClient(t *testing.T) {
	cookieStr := "phpdisk_info=UG0EMg1pUWkCMAVsWwhTPVQNBDdbMQBgBTdXMFVrBzNYZVBlBmQAO1dcB28AalRrVWQHNFpiXGlSNwU0UmAANlA3BGINOlFsAjQFNltlUz5UYgRgWzIAMgVjV2JVZgdnWGVQYAZhADxXNwdeAGtUPVUyBzRaNFwzUmcFbFJjADJQYw%3D%3D"
	c := NewMyClient(false)
	resp, err := c.Request("https://up.woozooo.com/mydisk.php", "get", cookieStr, nil, nil, nil)
	if nil != err {
		t.Error(err)
	}
	b, e := ioutil.ReadAll(resp.Body)
	if nil != e {
		t.Error(e)
	}
	t.Log(string(b), resp.Status, resp.Header.Get("Location"))
}

func TestMyClient_Request(t *testing.T) {
	resp, err := NewMyClient(false).Request("http://localhost:8888/t.php", "GET", "", nil, nil, nil)
	if nil != err {
		t.Log(err.Error())
	} else {
		b, _ := ioutil.ReadAll(resp.Body)
		t.Log(string(b))
	}
	for k, v := range resp.Header {
		t.Log(k, v)
	}
}
