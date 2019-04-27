package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/sclevine/agouti"
	"github.com/zhanglianxin/vvvdj-dl_go/config"
	"time"
)

const (
	UA = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_4) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/12.1 Safari/605.1.15"
)

var (
	url          string
	conf         *config.Config
	chromeDriver *agouti.WebDriver
	ua           string
	musicNames   []interface{}
)

func init() {
	conf = config.Load("config.toml")
	flag.StringVar(&url, "url", "http://www.vvvdj.com/radio/587.html", "The radio url")
	flag.Parse()
	chromeDriver = agouti.ChromeDriver(agouti.ChromeOptions("args", conf.Driver.Options["args"]))
}

func main() {
	// chromeDriver := agouti.ChromeDriver(agouti.ChromeOptions("args", conf.Driver.Options["args"]))
	if err := chromeDriver.Start(); nil != err {
		panic(err)
	}
	defer chromeDriver.Stop()

	page, err := chromeDriver.NewPage(agouti.Browser("chrome"))
	if nil != err {
		panic(err)
	}
	defer page.CloseWindow()

	if err := page.Navigate(url); nil != err {
		panic(err)
	}

	page.RunScript("return navigator.userAgent", nil, &ua) // get userAgent
	var ids string
	page.RunScript("return MUSICID;", nil, &ids) // get userAgent
	musicNames = getMusicNames(ids)
	for i := len(musicNames); i > 0; i-- {
		// TODO get element
	}
	fmt.Printf("%#v", musicNames)
	time.Sleep(10 * time.Second)
}

func getMusicNames(ids string) (data []interface{}) {
	url := "http://www.vvvdj.com/play/ajax/temp"
	headers := map[string]string{"X-Requested-With": "XMLHttpRequest", "User-Agent": ua}
	params := map[string]string{"ids": ids}
	hd := HtmlDocument{url: url, method: "GET", headers: headers, params: params}
	content := hd.getContent()
	var v interface{}
	if err := json.Unmarshal(content, &v); nil != err {
		panic(err)
	}
	var vv map[string]interface{}
	if err := json.Unmarshal([]byte(v.(string)), &vv); nil != err {
		panic(err)
	}
	if float64(200) == vv["Result"] {
		data = vv["Data"].([]interface{})
	}
	return data
}
