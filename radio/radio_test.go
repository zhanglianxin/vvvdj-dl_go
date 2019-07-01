package radio

import (
	"testing"
)

var (
	r *Radio
)

func init() {
	r = NewRadio("3454")
}

func TestRadio_GetPlayUrls(t *testing.T) {
	// r.getJsVarsViaOttoService()
	t.Logf("%#v", r.GetPlayUrls())
}
