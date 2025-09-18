package shared

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/samber/lo"
	"github.com/vegidio/go-sak/fetch"
	"github.com/vegidio/umd"
)

var cancelDownloads func()

func DownloadAll(
	media []umd.Media,
	directory string,
	parallel int,
	cookies []fetch.Cookie,
) <-chan *fetch.Response {
	cookieHeader := fetch.CookiesToHeader(cookies)

	f := fetch.New(map[string]string{
		"Cookie": cookieHeader,
	}, 10)

	requests := lo.Map(media, func(m umd.Media, _ int) *fetch.Request {
		filePath := CreateFilePath(directory, m)
		request, _ := f.NewRequest(m.Url, filePath)
		return request
	})

	resp, cancel := f.DownloadFiles(requests, parallel)
	cancelDownloads = cancel
	return resp
}

func CancelDownloads() {
	if cancelDownloads != nil {
		cancelDownloads()
		cancelDownloads = nil
	}
}

func ResponseToDownload(response *fetch.Response) Download {
	err := response.Error()
	bytes, _ := response.Bytes()

	return Download{
		Url:       response.Request.Url,
		FilePath:  response.Request.FilePath,
		Error:     err,
		IsSuccess: err == nil,
		Hash:      CreateFileHash(bytes),
	}
}

func CreateFilePath(directory string, media umd.Media) string {
	var t time.Time
	var err error

	n := media.Metadata["name"].(string)
	suffix := CreateHashSuffix(media.Url)

	// If the array of Media is coming from the JS code, the values in the Metadata map are strings
	timeStr, ok := media.Metadata["created"].(string)
	if ok {
		t, err = time.Parse(time.RFC3339, timeStr)
		if err != nil {
			t = time.Now()
		}
	} else {
		t = media.Metadata["created"].(time.Time)
	}

	timestamp := CreateTimestamp(t.Unix())

	ext := media.Extension
	if ext == "" {
		ext = "unknown"
	}

	fileName := fmt.Sprintf("%s-%s-%s.%s", n, timestamp, suffix, ext)
	return filepath.Join(directory, fileName)
}
