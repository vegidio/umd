package tests

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vegidio/umd"
	"github.com/vegidio/umd/internal/types"
)

func TestReddit_QuerySubreddit(t *testing.T) {
	if os.Getenv("GITHUB_ACTIONS") == "true" {
		t.Skip("This test doesn't work when executed from GitHub Actions")
	}

	const NumberOfPosts = 50

	extractor, _ := umd.New(nil).FindExtractor("https://www.reddit.com/r/PristineGirls/")
	resp, _ := extractor.QueryMedia(NumberOfPosts, nil, true)
	err := resp.Error()

	assert.NoError(t, err)
	assert.Equal(t, NumberOfPosts, len(resp.Media))
	assert.Equal(t, "subreddit", resp.Media[0].Metadata["source"])
	assert.Equal(t, "PristineGirls", resp.Media[0].Metadata["name"])
}

func TestReddit_QuerySubmissions(t *testing.T) {
	if os.Getenv("GITHUB_ACTIONS") == "true" {
		t.Skip("This test doesn't work when executed from GitHub Actions")
	}

	const NumberOfPosts = 50

	extractor, _ := umd.New(nil).FindExtractor("https://www.reddit.com/user/atomicbrunette18/")
	resp, _ := extractor.QueryMedia(NumberOfPosts, nil, true)
	err := resp.Error()

	assert.NoError(t, err)
	assert.Equal(t, NumberOfPosts, len(resp.Media))
	assert.Equal(t, "user", resp.Media[0].Metadata["source"])
	assert.Equal(t, "atomicbrunette18", resp.Media[0].Metadata["name"])
}

func TestReddit_QuerySingleSubmission(t *testing.T) {
	if os.Getenv("GITHUB_ACTIONS") == "true" {
		t.Skip("This test doesn't work when executed from GitHub Actions")
	}

	extractor, _ := umd.New(nil).FindExtractor("https://www.reddit.com/r/needysluts/comments/1aenk3e/if_im_wearing_this_for_our_date_you_have_bo/")
	resp, _ := extractor.QueryMedia(99999, nil, true)
	err := resp.Error()

	assert.NoError(t, err)
	assert.Equal(t, 1, len(resp.Media))
	assert.Equal(t, types.Video, resp.Media[0].Type)
	assert.Equal(t, "submission", resp.Media[0].Metadata["source"])
	assert.Equal(t, "needysluts", resp.Media[0].Metadata["name"])
}
