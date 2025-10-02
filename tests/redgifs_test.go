package tests

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	log "github.com/sirupsen/logrus"

	"github.com/stretchr/testify/assert"
	"github.com/vegidio/go-sak/fetch"
	"github.com/vegidio/umd"
)

func TestRedGifs_DownloadVideo(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "video.mp4")
	defer os.RemoveAll(tmpDir)

	extractor, _ := umd.New().FindExtractor("https://www.redgifs.com/watch/sturdycuddlyicefish")
	resp, _ := extractor.QueryMedia(99999, nil, true)
	<-resp.Done

	media := resp.Media[0]
	f := fetch.New(nil, 0)
	request, _ := f.NewRequest(media.Url, filePath)
	downloadResponse := f.DownloadFile(request)

	assert.NoError(t, downloadResponse.Error())
	assert.Equal(t, int64(15_212_770), downloadResponse.Size)
	assert.Equal(t, "sturdycuddlyicefish", media.Metadata["id"])
	assert.Equal(t, "sonya_18yo", media.Metadata["name"])
}

func TestRedGifs_FetchUser(t *testing.T) {
	extractor, _ := umd.New().FindExtractor("https://www.redgifs.com/users/atomicbrunette18")
	resp, _ := extractor.QueryMedia(180, nil, true)
	err := resp.Error()

	assert.NoError(t, err)
	assert.Equal(t, 180, len(resp.Media))
	assert.Equal(t, "user", resp.Media[0].Metadata["source"])
	assert.Equal(t, "atomicbrunette18", resp.Media[0].Metadata["name"])
}

func TestRedGifs_ReuseToken(t *testing.T) {
	// Create a buffer and redirect global log output to it
	var buf bytes.Buffer
	log.SetOutput(&buf)
	log.SetLevel(log.DebugLevel)

	// First query
	u := umd.New()
	extractor, _ := u.FindExtractor("https://www.redgifs.com/watch/sturdycuddlyicefish")
	r1, _ := extractor.QueryMedia(99999, nil, true)
	<-r1.Done

	// Second query
	u = umd.New().WithMetadata(r1.Metadata)
	extractor, _ = u.FindExtractor("https://www.redgifs.com/watch/ecstaticthickasiansmallclawedotter")
	r2, _ := extractor.QueryMedia(99999, nil, true)
	<-r2.Done

	// Check the log output
	output := buf.String()
	assert.Equal(t, 1, strings.Count(output, "Issuing new RedGifs token"))
	assert.Equal(t, 1, strings.Count(output, "Reusing RedGifs token"))
}
