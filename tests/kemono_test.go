package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vegidio/umd"
	"github.com/vegidio/umd/internal/types"
)

func TestKemono_QueryUser(t *testing.T) {
	const NumberOfPosts = 50

	extractor, _ := umd.New().FindExtractor("https://kemono.cr/patreon/user/4626321")
	resp, _ := extractor.QueryMedia(NumberOfPosts, nil, true)
	err := resp.Error()

	assert.NoError(t, err)
	assert.Equal(t, NumberOfPosts, len(resp.Media))
	assert.Equal(t, "patreon", resp.Media[0].Metadata["source"])
	assert.Equal(t, "haneame", resp.Media[0].Metadata["name"])
}

func TestKemono_QueryPost(t *testing.T) {
	extractor, _ := umd.New().FindExtractor("https://kemono.cr/patreon/user/4626321/post/122592054")
	resp, _ := extractor.QueryMedia(99999, nil, true)
	err := resp.Error()

	assert.NoError(t, err)
	assert.Equal(t, 4, len(resp.Media))
	assert.Equal(t, types.Image, resp.Media[0].Type)
	assert.Equal(t, "patreon", resp.Media[0].Metadata["source"])
	assert.Equal(t, "haneame", resp.Media[0].Metadata["name"])
}
