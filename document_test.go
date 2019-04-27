package main

import (
	"encoding/json"
	"io/ioutil"
	"testing"
)

func TestMakeRequest(t *testing.T) {
}

func TestResolveResult(t *testing.T) {
	b, _ := ioutil.ReadFile("test.txt")
	var v map[string]interface{}
	json.Unmarshal(b, &v)
	t.Logf("%#v", v["Result"])
	t.Logf("%#v", v["Data"])
}
