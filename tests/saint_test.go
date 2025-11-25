package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vegidio/umd"
)

func TestSaint_QueryVideo1(t *testing.T) {
	const NumberOfPosts = 1

	extractor, _ := umd.New().FindExtractor("https://saint2.su/embed/P9kEUyTHgJd")
	resp, _ := extractor.QueryMedia(99999, nil, true)
	err := resp.Error()

	assert.NoError(t, err)
	assert.Equal(t, NumberOfPosts, len(resp.Media))
	assert.Equal(t, "https://data.saint2.cr/data/P9kEUyTHgJd.mp4", resp.Media[0].Url)
	assert.Equal(t, "video", resp.Media[0].Metadata["source"])
	assert.Equal(t, "P9kEUyTHgJd", resp.Media[0].Metadata["name"])
}

func TestSaint_QueryVideo2(t *testing.T) {
	const NumberOfPosts = 1

	extractor, _ := umd.New().FindExtractor("https://saint2.cr/embed/wgqk6fjXugA")
	resp, _ := extractor.QueryMedia(99999, nil, true)
	err := resp.Error()

	assert.NoError(t, err)
	assert.Equal(t, NumberOfPosts, len(resp.Media))
	assert.Equal(t, "https://data.saint2.cr/data/wgqk6fjXugA.mp4", resp.Media[0].Url)
	assert.Equal(t, "video", resp.Media[0].Metadata["source"])
	assert.Equal(t, "wgqk6fjXugA", resp.Media[0].Metadata["name"])
}
