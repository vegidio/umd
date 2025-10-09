package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vegidio/umd"
)

func TestCyberdrop_QueryImage(t *testing.T) {
	const NumberOfPosts = 1

	extractor, _ := umd.New().FindExtractor("https://cyberdrop.me/f/YDHyWOicZvPKf")
	resp, _ := extractor.QueryMedia(99999, nil, true)
	err := resp.Error()

	assert.NoError(t, err)
	assert.Equal(t, NumberOfPosts, len(resp.Media))
	assert.Contains(t, resp.Media[0].Url, "https://k1-cd.cdn.gigachad-cdn.ru/api/file/d/YDHyWOicZvPKf")
	assert.Equal(t, "nanigazinski-20250902_232023-943519337-HmOXqZoc.jpg", resp.Media[0].Metadata["name"])
}

func TestCyberdrop_QueryAlbum(t *testing.T) {
	const NumberOfPosts = 2

	extractor, _ := umd.New().FindExtractor("https://cyberdrop.me/a/nU04Is4X")
	resp, _ := extractor.QueryMedia(99999, nil, true)
	err := resp.Error()

	assert.NoError(t, err)
	assert.Equal(t, NumberOfPosts, len(resp.Media))
	assert.Contains(t, resp.Media[0].Url, "https://k1-cd.cdn.gigachad-cdn.ru/api/file/d/nxeokHEPcugRd")
	assert.Equal(t, "nanigazinski-20250913_235925-249842057-YBy2RMbe.jpg", resp.Media[0].Metadata["name"])
}
