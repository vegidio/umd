package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vegidio/umd"
)

func TestErome_QueryAlbum1(t *testing.T) {
	const NumberOfMedia = 53

	extractor, _ := umd.New().FindExtractor("https://www.erome.com/a/YI93aUC3")
	resp, _ := extractor.QueryMedia(99999, nil, true)
	err := resp.Error()

	assert.NoError(t, err)
	assert.Equal(t, NumberOfMedia, len(resp.Media))
	assert.Equal(t, "YI93aUC3", resp.Media[0].Metadata["id"])
	assert.Equal(t, "likablewoman@Onlyfans", resp.Media[0].Metadata["name"])
}

func TestErome_QueryAlbum2(t *testing.T) {
	const NumberOfMedia = 21

	extractor, _ := umd.New().FindExtractor("https://www.erome.com/a/oidPGn1c")
	resp, _ := extractor.QueryMedia(99999, nil, true)
	err := resp.Error()

	assert.NoError(t, err)
	assert.Equal(t, NumberOfMedia, len(resp.Media))
	assert.Equal(t, "oidPGn1c", resp.Media[0].Metadata["id"])
	assert.Equal(t, "likablewoman", resp.Media[0].Metadata["name"])
}
