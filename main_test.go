package main

import (
	"encoding/json"
	"testing"
)

func TestGetMusicNames(t *testing.T) {
	ids := "105092,113138,113630,111317,115054,114892,115423,116077,117173,120361,119867,107345,126134,126042,126999,127489,127640,128832,129279,129497,130989,132255,137268,138469,138616,138574,139609,140554,140874,141264,128754,140097,139150,121960,142823,143568,143314,143990,144335,146460,146031,147281,147848,148005,143729,150007,150452,151908,152049,152436,152952,154209,154474,153785,159565,163987,158654,165189,166978,167171,168728,168896,170032"
	url := "http://www.vvvdj.com/play/ajax/temp"
	headers := map[string]string{"X-Requested-With": "XMLHttpRequest", "User-Agent": ua}
	params := map[string]string{"ids": ids}
	hd := HtmlDocument{url: url, method: "GET", headers: headers, params: params}
	content := hd.getContent()

	var i interface{}
	json.Unmarshal(content, &i)
	var v map[string]interface{}
	json.Unmarshal([]byte(i.(string)), &v)

	t.Logf("%#v", i)
	t.Logf("%#v", v)
	t.Logf("%#v", v["Result"])
}
