package update

import (
	"testing"

	"github.com/starudream/go-lib/testx"
)

func TestGetLatestReleaseName(t *testing.T) {
	name, err := GetLatestReleaseName()
	testx.P(t, err, name)
}

func TestDownloadReleaseFile(t *testing.T) {
	_, err := DownloadReleaseFile("v1.9.1", "", "", false)
	testx.P(t, err)
}
