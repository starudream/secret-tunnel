package updatex

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/minio/selfupdate"

	"github.com/starudream/go-lib/log"

	"github.com/starudream/secret-tunnel/constant"
)

type OSArch string

const (
	LinuxAMD64  OSArch = "linux-amd64"
	LinuxARM64  OSArch = "linux-arm64"
	DarwinAMD64 OSArch = "darwin-amd64"
	DarwinARM64 OSArch = "darwin-arm64"
	WindowAMD64 OSArch = "windows-amd64"
	WindowARM64 OSArch = "windows-arm64"
)

func GetOSArch() OSArch {
	return OSArch(runtime.GOOS + "-" + runtime.GOARCH)
}

func GetDownloadURL(osArch OSArch, version string) string {
	return fmt.Sprintf(constant.GitHub+"/releases/download/%s/secret-tunnel-%s-client-%s.zip", version, osArch, version)
}

var mirrors = map[string]*strings.Replacer{
	"fastgit": strings.NewReplacer("https://github.com/", "https://download.fastgit.org/"),
	"ghproxy": strings.NewReplacer("https://github.com/", "https://ghproxy.com/https://github.com/"),
}

func Try(version string) {
	l := log.With().Str("span", "update").Logger()

	update := func() io.Reader {
		url := GetDownloadURL(GetOSArch(), version)
		for name, replacer := range mirrors {
			bs, err := download(l, replacer.Replace(url))
			if err != nil {
				l.Warn().Str("name", name).Msgf("download update file error: %v", err)
				continue
			}
			l.Info().Str("name", name).Msgf("download update file success")
			return bytes.NewReader(bs)
		}
		return nil
	}()

	if update == nil {
		return
	}

	err := selfupdate.Apply(update, selfupdate.Options{})
	if err != nil {
		l.Error().Msgf("apply new file error: %v", err)
	}
}

var client = &http.Client{Timeout: 30 * time.Second}

func download(l log.L, url string) ([]byte, error) {
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	l.Debug().Msgf("get download file success, size: %s", resp.Header.Get("Content-Length"))

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	reader, err := zip.NewReader(bytes.NewReader(raw), int64(len(raw)))
	if err != nil {
		return nil, err
	}

	if len(reader.File) != 1 {
		return nil, fmt.Errorf("the zip archive must contain exactly one file")
	}

	rc, err := reader.File[0].Open()
	if err != nil {
		return nil, err
	}

	return io.ReadAll(rc)
}
