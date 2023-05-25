package update

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"time"

	"github.com/minio/selfupdate"

	"github.com/starudream/go-lib/codec/json"
	"github.com/starudream/go-lib/errx"
)

const (
	owner     = "starudream"
	repo      = "secret-tunnel"
	checksums = "checksums.txt"

	windows = "windows"
)

var (
	releaseMirrors = []string{
		"https://gh.ddlc.top/",
		"https://ghproxy.net/",
		"https://ghproxy.com/",
	}

	c = &http.Client{Timeout: 10 * time.Minute}
)

func Apply(update io.Reader) error {
	return selfupdate.Apply(update, selfupdate.Options{})
}

type xLatestRelease struct {
	Name string `json:"name"`
}

func GetLatestReleaseName() (string, error) {
	body, err := xGet("https://api.github.com/repos/" + owner + "/" + repo + "/releases/latest")
	if err != nil {
		return "", err
	}
	data, err := json.UnmarshalTo[*xLatestRelease](body)
	if err != nil {
		return "", err
	}
	return data.Name, nil
}

func DownloadReleaseFile(name, os, arch string, server bool) (io.Reader, error) {
	for i := 0; i < len(releaseMirrors); i++ {
		file, err := downloadReleaseFile(releaseMirrors[i], name, os, arch, server)
		if err == nil {
			return file, nil
		}
	}
	return nil, fmt.Errorf("download release file failed")
}

func downloadReleaseFile(mirror, name, os, arch string, server bool) (io.Reader, error) {
	sums, err := xGet(mirror + "https://github.com/" + owner + "/" + repo + "/releases/download/" + name + "/" + checksums)
	if err != nil {
		return nil, err
	}

	if os == "" {
		os = runtime.GOOS
	}
	if arch == "" {
		arch = runtime.GOARCH
	}
	target := "client"
	if server {
		target = "server"
	}
	ext := "tar.gz"
	if os == windows {
		ext = "zip"
	}

	filename := repo + "-" + target + "_" + name + "_" + os + "_" + arch + "." + ext

	var sum []byte
	for _, line := range bytes.Split(sums, []byte("\n")) {
		if bytes.Contains(line, []byte(filename)) {
			sum = bytes.TrimSpace(bytes.TrimSuffix(line, []byte(filename)))
		}
	}
	if len(sum) == 0 {
		return nil, fmt.Errorf("file not exists")
	}

	sum, err = hex.DecodeString(string(sum))
	if err != nil {
		return nil, err
	}

	file, err := xGet(mirror + "https://github.com/" + owner + "/" + repo + "/releases/download/" + name + "/" + filename)
	if err != nil {
		return nil, err
	}

	h := sha256.New()
	h.Write(file)
	hash := h.Sum(nil)
	if !bytes.Equal(sum, hash) {
		return nil, fmt.Errorf("checksum mismatch, expected %x, got %x", sum, hash)
	}

	if os == windows {
		zr, re := zip.NewReader(bytes.NewReader(file), int64(len(file)))
		if re != nil {
			return nil, re
		}
		for _, f := range zr.File {
			if f.Name == repo+"-"+target {
				raw, oe := f.OpenRaw()
				return raw, oe
			}
		}
	} else {
		gr, re := gzip.NewReader(bytes.NewReader(file))
		if re != nil {
			return nil, re
		}
		defer gr.Close()
		tr := tar.NewReader(gr)
		for {
			header, ne := tr.Next()
			if errx.Is(ne, io.EOF) {
				break
			}
			if ne != nil {
				return nil, ne
			}
			if header.Name == repo+"-"+target {
				raw, rae := io.ReadAll(tr)
				if rae != nil {
					return nil, rae
				}
				return bytes.NewReader(raw), nil
			}
		}
	}

	return nil, fmt.Errorf("file not exists")
}

func xGet(url string) ([]byte, error) {
	resp, err := c.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code: %d, body: %s", resp.StatusCode, bs)
	}

	return bs, err
}
