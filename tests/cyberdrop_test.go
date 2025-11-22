package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vegidio/umd"
)

func TestCyberdrop_QueryImage1(t *testing.T) {
	const NumberOfPosts = 1

	extractor, _ := umd.New().FindExtractor("https://cyberdrop.cr/f/zdqsYhkj8pzPy")
	resp, _ := extractor.QueryMedia(99999, nil, true)
	err := resp.Error()

	assert.NoError(t, err)
	assert.Equal(t, NumberOfPosts, len(resp.Media))
	assert.Contains(t, resp.Media[0].Url, "https://k1-cd.cdn.gigachad-cdn.ru/api/file/d/zdqsYhkj8pzPy")
	assert.Equal(t, "Lena (Kirill Chernyavsky) 6-gDpAfGLv.jpeg", resp.Media[0].Metadata["name"])
}

func TestCyberdrop_QueryImage2(t *testing.T) {
	const NumberOfPosts = 1

	extractor, _ := umd.New().FindExtractor("https://cyberdrop.me/f/pkRhIqRqNVJY2")
	resp, _ := extractor.QueryMedia(99999, nil, true)
	err := resp.Error()

	assert.NoError(t, err)
	assert.Equal(t, NumberOfPosts, len(resp.Media))
	assert.Contains(t, resp.Media[0].Url, "https://k1-cd.cdn.gigachad-cdn.ru/api/file/d/pkRhIqRqNVJY2")
	assert.Equal(t, "Lena (Kirill Chernyavsky) 4-ViHs90QA.jpeg", resp.Media[0].Metadata["name"])
}

func TestCyberdrop_QueryAlbum(t *testing.T) {
	const NumberOfPosts = 9

	extractor, _ := umd.New().FindExtractor("https://cyberdrop.cr/a/dHZ8Ffjy")
	resp, _ := extractor.QueryMedia(99999, nil, true)
	err := resp.Error()

	assert.NoError(t, err)
	assert.Equal(t, NumberOfPosts, len(resp.Media))
	assert.Contains(t, resp.Media[0].Url, "https://k1-cd.cdn.gigachad-cdn.ru/api/file/d/g70961KpRsBRI")
	assert.Equal(t, "Lena (Kirill Chernyavsky) 9-2AvfCVmu.jpeg", resp.Media[0].Metadata["name"])
}
