package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vegidio/umd"
)

func TestJpgFish_QueryImage1(t *testing.T) {
	const NumberOfPosts = 1

	extractor, _ := umd.New().FindExtractor("https://jpg6.su/img/NHVlaVI")
	resp, _ := extractor.QueryMedia(99999, nil, true)
	err := resp.Error()

	assert.NoError(t, err)
	assert.Equal(t, NumberOfPosts, len(resp.Media))
	assert.Equal(t, "https://simp6.selti-delivery.ru/images3/3648x5472_0e9cbce2c2fcab497a943bd192d90da48ca51abd591452cd.jpg", resp.Media[0].Url)
	assert.Equal(t, "image", resp.Media[0].Metadata["source"])
	assert.Equal(t, "solidsnake", resp.Media[0].Metadata["name"])
	assert.Equal(t, "3648x5472 0e9cbce2c2fcab497a943bd192d90da4", resp.Media[0].Metadata["title"])
}

func TestJpgFish_QueryImage2(t *testing.T) {
	const NumberOfPosts = 1

	extractor, _ := umd.New().FindExtractor("https://jpg7.cr/img/YAjShXh")
	resp, _ := extractor.QueryMedia(99999, nil, true)
	err := resp.Error()

	assert.NoError(t, err)
	assert.Equal(t, NumberOfPosts, len(resp.Media))
	assert.Equal(t, "https://simp6.selti-delivery.ru/images3/IMG_20240714_180333acaf99ec7bd1115a.jpg", resp.Media[0].Url)
	assert.Equal(t, "image", resp.Media[0].Metadata["source"])
	assert.Equal(t, "x34", resp.Media[0].Metadata["name"])
	assert.Equal(t, "IMG 20240714 180333", resp.Media[0].Metadata["title"])
}

func TestJpgFish_QueryImage_LongUrl(t *testing.T) {
	const NumberOfPosts = 1

	extractor, _ := umd.New().FindExtractor("https://jpg5.su/img/3648x5472-0e9cbce2c2fcab497a943bd192d90da4.NHVlaVI")
	resp, _ := extractor.QueryMedia(99999, nil, true)
	err := resp.Error()

	assert.NoError(t, err)
	assert.Equal(t, NumberOfPosts, len(resp.Media))
	assert.Equal(t, "https://simp6.selti-delivery.ru/images3/3648x5472_0e9cbce2c2fcab497a943bd192d90da48ca51abd591452cd.jpg", resp.Media[0].Url)
	assert.Equal(t, "image", resp.Media[0].Metadata["source"])
	assert.Equal(t, "solidsnake", resp.Media[0].Metadata["name"])
	assert.Equal(t, "3648x5472 0e9cbce2c2fcab497a943bd192d90da4", resp.Media[0].Metadata["title"])
}
