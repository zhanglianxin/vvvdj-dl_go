package main

import (
	"flag"
	"fmt"
	"github.com/zhanglianxin/vvvdj-dl_go/config"
	"runtime"
	"time"
)

var (
	start    time.Time
	radioUrl string
	tmpDir   = "tmp"
)

func init() {
	start = time.Now()

	config.Load("./config.toml")
	config.SetLog(start)
	runtime.GOMAXPROCS(runtime.NumCPU())

	flag.StringVar(&radioUrl, "radioUrl", "", "The radio url")
	flag.Parse()
}

func main() {
	dataDir := config.Conf.App.Data.Dir

	if "" == radioUrl {
		fmt.Println("params error")
		return
	}

	fmt.Println(dataDir)
}
