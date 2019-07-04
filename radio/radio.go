package radio

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/sirupsen/logrus"
	"github.com/zhanglianxin/vvvdj-dl_go/my-request"
	"io"
	"strings"
	"sync"
)

type Radio struct {
	host       string // host
	radioId    string // radio id
	musicUrl   string // current music url
	playUrls   map[string]string
	musicNames map[string]string
	playingId  string            // current music id
	musicIds   string            // music ids
	apsvr      string            // t.h
	playUrl    string            // current music src
	m4a        string            // "http://" + apsvr + ".vvvdj.com/face/" + file + ".mp4"
	source     *goquery.Document // current music page source
}

func NewRadio(radioId string) *Radio {
	return &Radio{
		radioId:  radioId,
		host:     "http://www.vvvdj.com",
		musicUrl: fmt.Sprintf("http://www.vvvdj.com/radio/%s.html", radioId),
		apsvr:    "t.h",
	}
}

func (r *Radio) GetPlayUrls() map[string]string {
	if "" == r.musicIds {
		r.getJsVarsViaOttoService()
	}
	r.getMusicNames()
	musicIds := strings.Split(r.musicIds, ",")
	length := len(musicIds)
	playUrls := make(map[string]string, length)

	var wg sync.WaitGroup
	wg.Add(length)
	for _, musicId := range musicIds {
		go func(musicId string) {
			defer wg.Done()
			radio := NewRadio(r.radioId)
			radio.musicUrl += fmt.Sprintf("?musicid=%s", musicId)
			radio.getM4a()
			playUrls[musicId] = radio.m4a
			if r.playingId == musicId {
				r.m4a = radio.m4a
			}
		}(musicId)
	}
	wg.Wait()

	r.playUrls = playUrls

	return r.playUrls
}

func (r *Radio) getMusicNames() {
	c := my_request.NewMyClient(false, 30)
	urlStr := "http://www.vvvdj.com/play/ajax/temp"
	headers := map[string]string{
		"X-Requested-With": "XMLHttpRequest",
	}
	params := map[string]string{
		"ids": r.musicIds,
	}
	res, err := c.Request(urlStr, "GET", "", nil, headers, params)
	if nil != err {
		logrus.Errorf("get music names: %s", err.Error())
		return
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	var v string
	if err := decoder.Decode(&v); nil != err {
		logrus.Errorf("decode music names 1 return error: %s", err.Error())
		return
	}

	decoder = json.NewDecoder(strings.NewReader(v))
	var vv map[string]interface{}
	if err := decoder.Decode(&vv); nil != err {
		logrus.Errorf("decode music names 2 return error: %s", err.Error())
		return
	}
	if float64(200) == vv["Result"] {
		slices := vv["Data"].([]interface{})
		musicNames := make(map[string]string, len(slices))
		for i := range slices {
			musicName := slices[i].(map[string]interface{})
			if musicId, ok := musicName["id"]; ok {
				musicNames[musicId.(string)] = musicName["musicname"].(string)
			}
		}
		r.musicNames = musicNames
	}
}

func (r *Radio) getM4a() {
	if "" == r.playUrl {
		r.getJsVarsViaOttoService()
	}

	r.m4a = fmt.Sprintf("http://%s.vvvdj.com/face/%s.mp4", r.apsvr, r.playUrl)
}

// get javascript variables
func (r *Radio) getJsVarsViaOttoService() {
	if nil == r.source {
		r.getSource()
	}
	infoScript := r.source.Find("div.radio_box + script").First().Text()
	if "" != infoScript {
		result, _ := removeLines(strings.NewReader(infoScript), []int{8, 9, 13})
		c := my_request.NewMyClient(false, 30)
		ottoServiceHost := "http://do:9998" // https://github.com/zhanglianxin/otto-service
		payload := strings.NewReader(result)
		res, err := c.Request(ottoServiceHost, "POST", "", payload, nil, nil)
		if nil != err {
			logrus.Errorf("post otto-service error: %s", err.Error())
			panic(err)
		}
		defer res.Body.Close()

		decoder := json.NewDecoder(res.Body)
		var m map[string]interface{}
		if err := decoder.Decode(&m); nil != err {
			logrus.Errorf("decode otto-service return error: %s", err.Error())
		}
		if val, ok := m["PLAYINGID"]; ok {
			r.playingId = fmt.Sprintf("%v", val)
		}
		if val, ok := m["MUSICID"]; ok {
			r.musicIds = val.(string)
		}
		if val, ok := m["playurl"]; ok {
			r.playUrl = val.(string)
		}
	}
}

// get radio page source
func (r *Radio) getSource() {
	c := my_request.NewMyClient(false, 30)
	urlStr := r.musicUrl
	headers := map[string]string{
		"Referer": r.host,
	}
	res, err := c.Request(urlStr, "GET", "", nil, headers, nil)
	defer res.Body.Close()

	if nil == err {
		doc, err := goquery.NewDocumentFromReader(res.Body)
		if nil != err {
			logrus.Errorf("new doc error: %s", err.Error())
			panic(err)
		}
		r.source = doc
	}
}

// https://stackoverflow.com/a/30708912
func readLine(r io.Reader, lineNum int) (line string, lastLine int, err error) {
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		lastLine++
		if lastLine == lineNum {
			// you can return sc.Bytes() if you need output in []bytes
			return sc.Text(), lastLine, sc.Err()
		}
	}
	return line, lastLine, io.EOF
}

func removeLines(r io.Reader, lineNums []int) (result string, lastLine int) {
	sc := bufio.NewScanner(r)

	for sc.Scan() {
		lastLine++
		if contains(lineNums, lastLine) {
			continue
		}
		result += fmt.Sprintln(sc.Text())
	}
	return result, lastLine
}

// https://stackoverflow.com/a/10485970
func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
