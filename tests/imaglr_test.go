package tests

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vegidio/go-sak/fetch"
	"github.com/vegidio/umd"
)

func TestImaglr_DownloadVideo(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "video.mp4")
	defer os.RemoveAll(tmpDir)

	extractor, _ := umd.New().FindExtractor("https://imaglr.com/post/5778297")
	resp, _ := extractor.QueryMedia(99999, nil, true)
	<-resp.Done

	media := resp.Media[0]
	f := fetch.New(nil, 0)
	request, _ := f.NewRequest(media.Url, filePath)
	downloadResponse := f.DownloadFile(request)

	assert.NoError(t, downloadResponse.Error())
	assert.Equal(t, int64(75_520_497), downloadResponse.Size)
}
