package shared

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/vegidio/go-sak/fetch"
	"github.com/vegidio/umd"
)

// cancelDownloads holds the cancel function for the current download session.
// Only one download session is supported at a time; starting a new one replaces the previous cancel function.
var cancelDownloads func()

// DownloadAll starts downloading every media item it can. Items whose URL can't even produce a request
// - a relative link, say - are returned separately as already-failed downloads, so that callers keep
// counting them and can put them in the report instead of losing them silently.
func DownloadAll(
	media []umd.Media,
	directory string,
	parallel int,
) (<-chan *fetch.Response, []Download) {
	f := fetch.New(nil, 10, false)

	requests := make([]*fetch.Request, 0, len(media))
	rejected := make([]Download, 0)

	for _, m := range media {
		filePath := CreateFilePath(directory, m)

		request, err := f.NewRequest(m.Url, filePath, m.Headers)
		if err != nil {
			rejected = append(rejected, Download{Url: m.Url, FilePath: filePath, Error: err})
			continue
		}

		requests = append(requests, request)
	}

	resp, cancel := f.DownloadFiles(requests, parallel)
	cancelDownloads = cancel
	return resp, rejected
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
