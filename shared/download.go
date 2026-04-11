package shared

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/samber/lo"
	"github.com/vegidio/go-sak/fetch"
	"github.com/vegidio/umd"
)

// cancelDownloads holds the cancel function for the current download session.
// Only one download session is supported at a time; starting a new one replaces the previous cancel function.
var cancelDownloads func()

func DownloadAll(
	media []umd.Media,
	directory string,
	parallel int,
) <-chan *fetch.Response {
	f := fetch.New(nil, 10, false)

	requests := lo.Map(media, func(m umd.Media, _ int) *fetch.Request {
		filePath := CreateFilePath(directory, m)
		request, _ := f.NewRequest(m.Url, filePath, m.Headers)
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

	return Download{
		Url:       response.Request.Url,
		FilePath:  response.Request.FilePath,
		Error:     err,
		IsSuccess: err == nil,
		Hash:      response.Hash,
	}
}

func CreateFilePath(directory string, media umd.Media) string {
	var t time.Time

	n, _ := media.Metadata["name"].(string)
	if n == "" {
		n = "unknown"
	}
	suffix := CreateHashSuffix(media.Url)

	// If the array of Media is coming from the JS code, the values in the Metadata map are strings
	switch v := media.Metadata["created"].(type) {
	case string:
		var err error
		t, err = time.Parse(time.RFC3339, v)
		if err != nil {
			t = time.Now()
		}
	case time.Time:
		t = v
	default:
		t = time.Now()
	}

	timestamp := CreateTimestamp(t.Unix())

	ext := media.Extension
	if ext == "" {
		ext = "unknown"
	}

	fileName := fmt.Sprintf("%s-%s-%s.%s", n, timestamp, suffix, ext)
	return filepath.Join(directory, fileName)
}
