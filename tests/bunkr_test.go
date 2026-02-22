package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vegidio/umd"
)

func TestBunkr(t *testing.T) {
	tests := []struct {
		name           string
		url            string
		numberOfMedia  int
		firstMediaName string
	}{
		{
			name:           "QueryImage",
			url:            "https://bunkr.cr/f/1440x2560_89acb8b819c3c63089281462663cc0bb-AtmU95Y9.jpg",
			numberOfMedia:  1,
			firstMediaName: "1440x2560_89acb8b819c3c63089281462663cc0bb.jpg",
		},
		{
			name:           "QueryVideo",
			url:            "https://bunkr.cr/v/0h44nakg3gkt0wh1nb05g_source-ARUJIdg7.mp4",
			numberOfMedia:  1,
			firstMediaName: "0h44nakg3gkt0wh1nb05g_source.mp4",
		},
		{
			name:           "QueryAlbum",
			url:            "https://bunkr.cr/a/v40v0xW1",
			numberOfMedia:  31,
			firstMediaName: "1920x2560_da96fee1aa635b45a510a3d260b5c678.jpg",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			extractor, _ := umd.New().FindExtractor(tt.url)
			resp, _ := extractor.QueryMedia(99999, nil, true)
			err := resp.Error()

			assert.NoError(t, err)
			assert.Equal(t, tt.numberOfMedia, len(resp.Media))
			assert.Equal(t, tt.firstMediaName, resp.Media[0].Metadata["name"])
		})
	}
}
