package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/sclevine/agouti"
	"github.com/sirupsen/logrus"
	"github.com/zhanglianxin/vvvdj-dl_go/config"
	"io"
	"math/rand"
	"os"
	"strings"
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
	logName := time.Now().Format("2006-01-02") + ".log"
	file, err := os.OpenFile(logName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if nil != err {
		panic(err)
	}
	logrus.SetOutput(file)
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
		currentUrl, _ := page.URL()
		audio := page.Find("audio#jp_audio_0[src]")
		anchorNext := page.FindByXPath("//*[@id=\"ico-next\"]/a")
		src, _ := audio.Attribute("src")
		ch := make(chan string)
		go func() {
			for "" == src {
				time.Sleep(250 * time.Millisecond)
				src, _ = audio.Attribute("src")
			}
			ch <- src
		}()
		select {
		case <-ch:
			sl := strings.SplitAfter(strings.Split(src, "?")[0], "/")
			fileName := sl[len(sl)-1]
			fmt.Println("fileName:", fileName)
			save2File(src, fileName, "data/")
		case <-time.After(5 * time.Second):
			panic("get src attribute timeout")
		}

		logrus.Infof("count: %d, current: %s, src: %s", i, currentUrl, src)
		if b, _ := anchorNext.Visible(); b {
			anchorNext.Click()
			url, _ := page.URL()
			ch := make(chan int)
			go func() {
				for currentUrl == url {
					time.Sleep(250 * time.Microsecond)
					url, _ = page.URL()
				}
				ch <- 1
			}()
			select {
			case <-ch:
			case <-time.After(5 * time.Second):
				panic("not forward")
			}
		}
		time.Sleep(time.Duration(randNum(1000, 3000)) * time.Microsecond)
	}
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

func save2File(url string, name string, path string) {
	if "" == path {
		path = "data/"
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 0755)
	}
	headers := map[string]string{"User-Agent": ua}
	params := map[string]string{}
	sd := StreamDocument{url: url, method: "GET", headers: headers, params: params}
	content := sd.getContent()
	defer content.Close()

	out, err := os.Create(path + name)
	defer out.Close()
	if nil != err {
		panic(err)
	}

	_, err = io.Copy(out, content)
	if nil != err {
		panic(err)
	}
}

func randNum(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}
