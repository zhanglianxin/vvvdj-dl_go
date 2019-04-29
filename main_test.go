package main

import (
	"encoding/json"
	"fmt"
	"github.com/sclevine/agouti"
	"math/rand"
	"strings"
	"testing"
	"time"
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

func TestSlice(t *testing.T) {
	s := "ttp://t.h.vvvdj.com/face/c2/2014/11/105092-e9c425.mp4?upt=038b56911587969253&play.mp4"
	sl := strings.SplitAfter(strings.Split(s, "?")[0], "/")
	s = sl[len(sl)-1]
	t.Logf("%#v", s)
}

func TestElemProp(t *testing.T) {
	if err := chromeDriver.Start(); nil != err {
		panic(err)
	}
	defer chromeDriver.Stop()

	page, err := chromeDriver.NewPage(agouti.Browser("chrome"))
	if nil != err {
		panic(err)
	}
	defer page.CloseWindow()

	url := "https://blog.test/"
	if err := page.Navigate(url); nil != err {
		panic(err)
	}

	elem := page.FindByXPath("/html/body/div/div/div[2]/a[1]")
	src, _ := elem.Attribute("src") // ""
	href, _ := elem.Attribute("href")
	t.Logf("%#v\n%#v", src, href)
}

func TestChan(t *testing.T) {
	ch := make(chan int)
	go func() {
		time.Sleep(5 * time.Second)
		ch <- 1
	}()
	a := <-ch
	fmt.Println(a)
}

func TestChan1(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	fmt.Println(rand.Intn(30))
	ch := make(chan string)
	fmt.Println(ch)
	var src string
	go func() {
		for {
			time.Sleep(1 * time.Second)
			src = "a"
			ch <- src
		}
	}()

	select {
	case <-ch:
		fmt.Println(ch, <-ch, src, "aha")
	case <-time.After(2 * time.Second):
		fmt.Println(ch, "oppos")
	}
}

func TestSave2File(t *testing.T) {
	mocks := []struct {
		params   map[string]string
		expected string
	}{
		{map[string]string{"url": "http://t.h.vvvdj.com/face/c2/2018/04/159717-7ec72a.mp4", "name": "159717-7ec72a.mp4", "path": "data/3454/"}, "f6221a402f39abd0e586d77073ca21af"},
	}

	for i := range mocks {
		save2File(mocks[i].params["url"], mocks[i].params["name"], mocks[i].params["path"])
		sum := md5Sum(mocks[i].params["path"]+mocks[i].params["name"])
		if mocks[i].expected != sum {
			t.Errorf("expected: [%v], \nactually: [%v]", mocks[i].expected, sum)
		}
	}
}
