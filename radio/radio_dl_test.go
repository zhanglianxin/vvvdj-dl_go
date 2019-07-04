package radio

import (
	"testing"
)

func TestRadioDl_Download(t *testing.T) {
	r = NewRadio("3454")
	r.GetPlayUrls()
	NewRadioDl().Download(r, "../data")
}
