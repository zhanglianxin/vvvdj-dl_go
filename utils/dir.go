package utils

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"runtime"
)

func CheckDir(path string) {
	canonical := "drwxr-xr-x"
	info, err := os.Stat(path)
	if nil != err {
		if !os.IsNotExist(err) {
			os.RemoveAll(path)
		}
		// Create
		err := os.MkdirAll(path, os.ModePerm)
		if nil != err {
			fmt.Println("wrong")
			logrus.Errorf("path: [%v], err: [%v]", path, err.Error())
		}
		info, _ = os.Stat(path)
	}
	if isWindows() {
		return // just return
	}
	mode := info.Mode().String()
	if canonical != mode {
		fmt.Println(mode)
		logrus.Errorf("path: [%v], mode: [%v]", path, mode)
	}
}

func isWindows() bool {
	return "windows" == runtime.GOOS
}
