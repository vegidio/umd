package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vegidio/umd"
)

func TestFapello_QueryPost(t *testing.T) {
	const NumberOfPosts = 1

	extractor, _ := umd.New(nil).FindExtractor("https://fapello.com/eva-padlock/1552/")
	resp, _ := extractor.QueryMedia(99999, nil, true)
	err := resp.Error()

	assert.NoError(t, err)
	assert.Equal(t, NumberOfPosts, len(resp.Media))
	assert.Equal(t, "post", resp.Media[0].Metadata["source"])
	assert.Equal(t, "eva-padlock", resp.Media[0].Metadata["name"])
}

func TestFapello_QueryModel(t *testing.T) {
	const NumberOfPosts = 98

	extractor, _ := umd.New(nil).FindExtractor("https://fapello.com/darja-sobakinskaja/")
	resp, _ := extractor.QueryMedia(99999, nil, true)
	err := resp.Error()

	assert.NoError(t, err)
	assert.Equal(t, NumberOfPosts, len(resp.Media))
	assert.Equal(t, "model", resp.Media[0].Metadata["source"])
	assert.Equal(t, "darja-sobakinskaja", resp.Media[0].Metadata["name"])
}
