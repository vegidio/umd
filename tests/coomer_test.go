package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vegidio/umd"
	"github.com/vegidio/umd/internal/types"
)

func TestCoomer_QueryUser(t *testing.T) {
	const NumberOfPosts = 50

	extractor, _ := umd.New().FindExtractor("https://coomer.st/onlyfans/user/melindalondon")
	resp, _ := extractor.QueryMedia(NumberOfPosts, nil, true)
	err := resp.Error()

	assert.NoError(t, err)
	assert.Equal(t, NumberOfPosts, len(resp.Media))
	assert.Equal(t, "onlyfans", resp.Media[0].Metadata["source"])
	assert.Equal(t, "melindalondon", resp.Media[0].Metadata["name"])
}

func TestCoomer_QueryPost(t *testing.T) {
	extractor, _ := umd.New().FindExtractor("https://coomer.st/onlyfans/user/melindalondon/post/357160243")
	resp, _ := extractor.QueryMedia(99999, nil, true)
	err := resp.Error()

	assert.NoError(t, err)
	assert.Equal(t, 1, len(resp.Media))
	assert.Equal(t, types.Image, resp.Media[0].Type)
	assert.Equal(t, "onlyfans", resp.Media[0].Metadata["source"])
	assert.Equal(t, "melindalondon", resp.Media[0].Metadata["name"])
}

func TestCoomer_QueryPostWithRevisions(t *testing.T) {
	extractor, _ := umd.New().FindExtractor("https://coomer.st/fansly/user/286621667281612800/post/290925242786787328")
	resp, _ := extractor.QueryMedia(99999, nil, true)
	err := resp.Error()

	assert.NoError(t, err)
	assert.Equal(t, 19, len(resp.Media))
	assert.Equal(t, types.Image, resp.Media[0].Type)
	assert.Equal(t, "fansly", resp.Media[0].Metadata["source"])
	assert.Equal(t, "Morgpie", resp.Media[0].Metadata["name"])
}
