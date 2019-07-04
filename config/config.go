package config

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/sirupsen/logrus"
	"github.com/zhanglianxin/vvvdj-dl_go/utils"
	"os"
	"time"
)

type Config struct {
	App *struct {
		Log *struct {
			Dir string `toml:"dir"`
		} `toml:"log"`
		Data *struct {
			Dir string `toml:"dir"`
		} `toml:"data"`
	} `toml:"app"`
}

var (
	Conf *Config
)

func Load(file string) *Config {
	if Conf == nil {
		if _, err := toml.DecodeFile(file, &Conf); nil != err {
			panic(err)
		}
	}
	return Conf
}

func SetLog(t time.Time) {
	if nil == Conf {
		panic("Please load Conf file first!!")
	}

	utils.CheckDir(Conf.App.Log.Dir)
	logName := fmt.Sprintf("%s/%s.log", Conf.App.Log.Dir, t.Format("2006-01-02"))
	file, err := os.OpenFile(logName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if nil != err {
		panic(err)
	}
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors: true,
	})
	logrus.SetOutput(file)
}
