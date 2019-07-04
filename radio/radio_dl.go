package radio

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/zhanglianxin/vvvdj-dl_go/my-request"
	"github.com/zhanglianxin/vvvdj-dl_go/utils"
	"io"
	"net/http"
	"os"
	"path"
	"sync"
)

type RadioDl struct {
}

func NewRadioDl() *RadioDl {
	return &RadioDl{}
}

func (dl *RadioDl) Download(r *Radio, dataDir string) {
	dir := dataDir + string(os.PathSeparator) + r.radioId + string(os.PathSeparator)
	utils.CheckDir(dir)

	m := r.playUrls
	var wg sync.WaitGroup
	wg.Add(len(m))
	for k := range m {
		go func(k string) {
			defer wg.Done()
			link := m[k]
			reader, err := getMediaContent(link)
			if nil != err {
				logrus.Errorf("get content error: %s", err.Error())
			} else {
				basename := k
				if musicName, ok := r.musicNames[k]; ok {
					basename = musicName
				}
				fname := dir + basename + path.Ext(link)
				written := save2File(reader, fname)
				if 0 == written {
					logrus.Errorf("save file error: %s", fname)
				}
			}
			reader.Close()
		}(k)
	}
	wg.Wait()
}

func getMediaContent(urlStr string) (io.ReadCloser, error) {
	c := my_request.NewMyClient(true, 120)
	res, err := c.Request(urlStr, "GET", "", nil, nil, nil)
	if nil != err {
		return nil, err
	}
	// defer res.Body.Close()
	if http.StatusOK == res.StatusCode {
		return res.Body, nil
	} else {
		return res.Body, errors.New(fmt.Sprintf("status code: %d", res.StatusCode))
	}
}

func save2File(content io.Reader, dst string) int {
	out, err := os.Create(dst)
	defer out.Close()
	if nil != err {
		panic(err)
	}
	written, err := io.Copy(out, content)
	if nil != err {
		logrus.Errorf("written: [%v], err: [%v]", written, err.Error())
	} else {
		logrus.Infof("file: [%v], written: [%d]", out.Name(), written)
	}
	return int(written)
}
