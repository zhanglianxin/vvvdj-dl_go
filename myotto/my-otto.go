package myotto

import (
	"github.com/robertkrimen/otto"
)

type MyOtto struct {
	otto.Otto
}

func New() *MyOtto {
	return &MyOtto{*otto.New()}
}

func (m MyOtto) GetString(name string) string {
	if v, e := m.Otto.Get(name); nil != e {
		return ""
	} else {
		if s, e := v.ToString(); nil != e {
			return ""
		} else {
			return s
		}
	}
}
