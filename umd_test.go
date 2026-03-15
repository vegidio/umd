package umd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vegidio/umd/internal/types"
)

// region - New

func TestNew_ReturnsNonNil(t *testing.T) {
	u := New()
	assert.NotNil(t, u)
}

func TestNew_InitializesMetadata(t *testing.T) {
	u := New()
	assert.NotNil(t, u.metadata)
}

// endregion

// region - WithMetadata

func TestWithMetadata_SetsMetadata(t *testing.T) {
	meta := Metadata{
		types.Reddit: {"token": "abc123"},
	}
	u := New().WithMetadata(meta)
	assert.Equal(t, "abc123", u.metadata[types.Reddit]["token"])
}

func TestWithMetadata_ReturnsSelf(t *testing.T) {
	u := New()
	result := u.WithMetadata(make(Metadata))
	assert.Same(t, u, result)
}

// endregion

// region - FindExtractor

func TestFindExtractor_Bunkr(t *testing.T) {
	ext, err := New().FindExtractor("https://bunkr.cr/a/test-album")
	assert.NoError(t, err)
	assert.Equal(t, types.Bunkr, ext.Type())
}

func TestFindExtractor_Coomer(t *testing.T) {
	ext, err := New().FindExtractor("https://coomer.su/onlyfans/user/test")
	assert.NoError(t, err)
	assert.Equal(t, types.Coomer, ext.Type())
}

func TestFindExtractor_Kemono(t *testing.T) {
	ext, err := New().FindExtractor("https://kemono.cr/patreon/user/test")
	assert.NoError(t, err)
	assert.Equal(t, types.Kemono, ext.Type())
}

func TestFindExtractor_Cyberdrop(t *testing.T) {
	ext, err := New().FindExtractor("https://cyberdrop.cr/a/test-album")
	assert.NoError(t, err)
	assert.Equal(t, types.Cyberdrop, ext.Type())
}

func TestFindExtractor_Erome(t *testing.T) {
	ext, err := New().FindExtractor("https://erome.com/a/test-album")
	assert.NoError(t, err)
	assert.Equal(t, types.Erome, ext.Type())
}

func TestFindExtractor_Fapello(t *testing.T) {
	ext, err := New().FindExtractor("https://fapello.com/test-model/")
	assert.NoError(t, err)
	assert.Equal(t, types.Fapello, ext.Type())
}

func TestFindExtractor_Reddit(t *testing.T) {
	ext, err := New().FindExtractor("https://reddit.com/r/test")
	assert.NoError(t, err)
	assert.Equal(t, types.Reddit, ext.Type())
}

func TestFindExtractor_RedGifs(t *testing.T) {
	ext, err := New().FindExtractor("https://redgifs.com/watch/test")
	assert.NoError(t, err)
	assert.Equal(t, types.RedGifs, ext.Type())
}

func TestFindExtractor_UnsupportedURL(t *testing.T) {
	ext, err := New().FindExtractor("https://unsupported-site.com/page")
	assert.Error(t, err)
	assert.Nil(t, ext)
	assert.Contains(t, err.Error(), "no extractor found")
}

func TestFindExtractor_SimpCity_RequiresCookies(t *testing.T) {
	// SimpCity requires cookies, so FindExtractor returns an error even though the extractor matches
	ext, err := New().FindExtractor("https://simpcity.cr/threads/test-thread")
	assert.Error(t, err)
	assert.Nil(t, ext)
	assert.Contains(t, err.Error(), "cookies")
}

// endregion
