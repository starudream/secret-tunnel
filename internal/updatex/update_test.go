package updatex

import (
	"testing"
)

func TestGetDownloadURL(t *testing.T) {
	t.Log(GetDownloadURL(GetOSArch(), "v1.4.0"))
}

func TestTry(t *testing.T) {
	Try("v1.4.0")
}
