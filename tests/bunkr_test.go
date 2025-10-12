package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vegidio/umd"
)

func TestBunkr_QueryImage(t *testing.T) {
	const NumberOfPosts = 1

	extractor, _ := umd.New().FindExtractor("https://bunkr.cr/f/1440x2560_89acb8b819c3c63089281462663cc0bb-AtmU95Y9.jpg")
	resp, _ := extractor.QueryMedia(99999, nil, true)
	err := resp.Error()

	assert.NoError(t, err)
	assert.Equal(t, NumberOfPosts, len(resp.Media))
	assert.Equal(t, "https://taquito.bunkr.ru/1440x2560_89acb8b819c3c63089281462663cc0bb-AtmU95Y9.jpg", resp.Media[0].Url)
	assert.Equal(t, "1440x2560_89acb8b819c3c63089281462663cc0bb.jpg", resp.Media[0].Metadata["name"])
}

func TestBunkr_QueryVideo(t *testing.T) {
	const NumberOfPosts = 1

	extractor, _ := umd.New().FindExtractor("https://bunkr.cr/f/0h6ztncb3m26odi997lvb_source-cYy4SjoG.mp4")
	resp, _ := extractor.QueryMedia(99999, nil, true)
	err := resp.Error()

	assert.NoError(t, err)
	assert.Equal(t, NumberOfPosts, len(resp.Media))
	assert.Equal(t, "https://taquito.bunkr.ru/0h6ztncb3m26odi997lvb_source-cYy4SjoG.mp4", resp.Media[0].Url)
	assert.Equal(t, "0h6ztncb3m26odi997lvb_source.mp4", resp.Media[0].Metadata["name"])
}

func TestBunkr_QueryAlbum(t *testing.T) {
	const NumberOfPosts = 32

	extractor, _ := umd.New().FindExtractor("https://bunkr.cr/a/v40v0xW1")
	resp, _ := extractor.QueryMedia(99999, nil, true)
	err := resp.Error()

	assert.NoError(t, err)
	assert.Equal(t, NumberOfPosts, len(resp.Media))
	assert.Equal(t, "https://taquito.bunkr.ru/1920x2560_da96fee1aa635b45a510a3d260b5c678-Wuazp1bG.jpg", resp.Media[0].Url)
	assert.Equal(t, "1920x2560_da96fee1aa635b45a510a3d260b5c678.jpg", resp.Media[0].Metadata["name"])
}
