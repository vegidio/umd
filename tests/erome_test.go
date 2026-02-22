package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vegidio/umd"
)

func TestErome_QueryAlbum(t *testing.T) {
	tests := []struct {
		name          string
		url           string
		numberOfMedia int
		albumName     string
		albumTitle    string
	}{
		{
			name:          "Album1",
			url:           "https://www.erome.com/a/YI93aUC3",
			numberOfMedia: 53,
			albumName:     "YI93aUC3",
			albumTitle:    "likablewoman@Onlyfans",
		},
		{
			name:          "Album2",
			url:           "https://www.erome.com/a/oidPGn1c",
			numberOfMedia: 21,
			albumName:     "oidPGn1c",
			albumTitle:    "likablewoman",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			extractor, _ := umd.New().FindExtractor(tt.url)
			resp, _ := extractor.QueryMedia(99999, nil, true)
			err := resp.Error()

			assert.NoError(t, err)
			assert.Equal(t, tt.numberOfMedia, len(resp.Media))
			assert.Equal(t, tt.albumName, resp.Media[0].Metadata["name"])
			assert.Equal(t, tt.albumTitle, resp.Media[0].Metadata["title"])
		})
	}
}
