package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vegidio/umd"
)

func TestSaint_QueryVideo(t *testing.T) {
	const NumberOfPosts = 1

	extractor, _ := umd.New(nil).FindExtractor("https://saint2.su/embed/P9kEUyTHgJd")
	resp, _ := extractor.QueryMedia(99999, nil, true)
	err := resp.Error()

	assert.NoError(t, err)
	assert.Equal(t, NumberOfPosts, len(resp.Media))
	assert.Equal(t, "https://data.saint2.cr/data/P9kEUyTHgJd.mp4", resp.Media[0].Url)
	assert.Equal(t, "video", resp.Media[0].Metadata["source"])
	assert.Equal(t, "P9kEUyTHgJd", resp.Media[0].Metadata["name"])
}
