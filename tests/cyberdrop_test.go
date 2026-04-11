package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vegidio/umd"
)

func TestCyberdrop(t *testing.T) {
	t.Run("QueryImage1", func(t *testing.T) {
		const NumberOfPosts = 1

		extractor, _ := umd.New().FindExtractor("https://cyberdrop.cr/f/YSG2LEVWPg7Ca")
		resp, _ := extractor.QueryMedia(99999, nil, true)
		err := resp.Error()

		assert.NoError(t, err)
		assert.Equal(t, NumberOfPosts, len(resp.Media))
		assert.Contains(t, resp.Media[0].Url, "https://k1-cd.cdn.gigachad-cdn.ru/api/file/d/YSG2LEVWPg7Ca")
		assert.Equal(t, "Irina Part 00198-Vop7PdK3.jpg", resp.Media[0].Metadata["name"])
	})

	t.Run("QueryImage2", func(t *testing.T) {
		const NumberOfPosts = 1

		extractor, _ := umd.New().FindExtractor("https://cyberdrop.cr/f/4ijoCSrnu6sml")
		resp, _ := extractor.QueryMedia(99999, nil, true)
		err := resp.Error()

		assert.NoError(t, err)
		assert.Equal(t, NumberOfPosts, len(resp.Media))
		assert.Contains(t, resp.Media[0].Url, "https://k1-cd.cdn.gigachad-cdn.ru/api/file/d/4ijoCSrnu6sml")
		assert.Equal(t, "Irina Part 00200-YPPqN4iL.jpg", resp.Media[0].Metadata["name"])
	})

	t.Run("QueryAlbum", func(t *testing.T) {
		const NumberOfPosts = 4

		extractor, _ := umd.New().FindExtractor("https://cyberdrop.cr/a/CUa8CjYU")
		resp, _ := extractor.QueryMedia(99999, nil, true)
		err := resp.Error()

		assert.NoError(t, err)
		assert.Equal(t, NumberOfPosts, len(resp.Media))
		assert.Contains(t, resp.Media[0].Url, "https://k1-cd.cdn.gigachad-cdn.ru/api/file/d/4ijoCSrnu6sml")
		assert.Equal(t, "Irina Part 00200-YPPqN4iL.jpg", resp.Media[0].Metadata["name"])
	})
}
