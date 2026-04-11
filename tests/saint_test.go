package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vegidio/umd"
)

func TestSaint(t *testing.T) {
	t.Run("QueryVideo1", func(t *testing.T) {
		const NumberOfPosts = 1

		extractor, _ := umd.New().FindExtractor("https://saint2.su/embed/P9kEUyTHgJd")
		resp, _ := extractor.QueryMedia(99999, nil, true)
		err := resp.Error()

		assert.NoError(t, err)
		assert.Contains(t, resp.Media[0].Url, "/data/P9kEUyTHgJd.mp4")
		assert.Equal(t, NumberOfPosts, len(resp.Media))
		assert.Equal(t, "video", resp.Media[0].Metadata["source"])
		assert.Equal(t, "P9kEUyTHgJd", resp.Media[0].Metadata["name"])
	})

	t.Run("QueryVideo2", func(t *testing.T) {
		const NumberOfPosts = 1

		extractor, _ := umd.New().FindExtractor("https://saint2.cr/embed/wgqk6fjXugA")
		resp, _ := extractor.QueryMedia(99999, nil, true)
		err := resp.Error()

		assert.NoError(t, err)
		assert.Contains(t, resp.Media[0].Url, "/data/wgqk6fjXugA.mp4")
		assert.Equal(t, NumberOfPosts, len(resp.Media))
		assert.Equal(t, "video", resp.Media[0].Metadata["source"])
		assert.Equal(t, "wgqk6fjXugA", resp.Media[0].Metadata["name"])
	})
}
