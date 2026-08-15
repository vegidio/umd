package shared

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vegidio/umd"
)

// Reddit hands out relative links such as "/r/subreddit" for crossposts. Those can never be
// downloaded, but they must not vanish either: the caller counts them and puts them in the report.
func TestDownloadAll_ReturnsMediaItThatCannotDownload(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("hello"))
	}))
	defer server.Close()

	media := []umd.Media{
		{Url: server.URL + "/a.jpg", Extension: "jpg"},
		{Url: "/r/some-subreddit"},
		{Url: "ftp://example.com/b.jpg", Extension: "jpg"},
	}

	result, rejected := DownloadAll(media, t.TempDir(), 2)

	require.Len(t, rejected, 2)
	assert.ElementsMatch(t,
		[]string{"/r/some-subreddit", "ftp://example.com/b.jpg"},
		[]string{rejected[0].Url, rejected[1].Url},
	)

	for _, d := range rejected {
		assert.False(t, d.IsSuccess)
		assert.Error(t, d.Error)
		assert.NotEmpty(t, d.FilePath)
	}

	// Only the downloadable item produces a response, so the caller's total is len(media) - len(rejected)
	responses := 0
	for range result {
		responses++
	}

	assert.Equal(t, 1, responses)
}

func TestDownloadAll_NoRejectsWhenEveryUrlIsUsable(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("hello"))
	}))
	defer server.Close()

	media := []umd.Media{
		{Url: server.URL + "/a.jpg", Extension: "jpg"},
		{Url: server.URL + "/b.jpg", Extension: "jpg"},
	}

	result, rejected := DownloadAll(media, t.TempDir(), 2)

	assert.Empty(t, rejected)

	responses := 0
	for r := range result {
		require.NoError(t, r.Error())
		responses++
	}

	assert.Equal(t, len(media), responses)
}
