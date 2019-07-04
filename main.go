package main

import (
	"flag"
	"fmt"
	"github.com/zhanglianxin/vvvdj-dl_go/config"
	"github.com/zhanglianxin/vvvdj-dl_go/radio"
	"runtime"
	"time"
)

var (
	start   time.Time
	radioId string
)

func init() {
	start = time.Now()

	config.Load("./config.toml")
	config.SetLog(start)
	runtime.GOMAXPROCS(runtime.NumCPU())

	flag.StringVar(&radioId, "radioId", "", "The radio id")
	flag.Parse()
}

func main() {
	dataDir := config.Conf.App.Data.Dir
	if "" == radioId {
		fmt.Println("params error")
		return
	}

	r := radio.NewRadio(radioId)
	r.GetPlayUrls()
	rdl := radio.NewRadioDl()
	rdl.Download(r, dataDir)
}
