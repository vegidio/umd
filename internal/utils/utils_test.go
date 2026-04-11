package utils

import (
	"testing"
	"time"

	saktypes "github.com/vegidio/go-sak/types"
	"github.com/vegidio/umd/internal/types"

	"github.com/stretchr/testify/assert"
)

// region - HasHost

func TestHasHost_ExactMatch(t *testing.T) {
	assert.True(t, HasHost("https://example.com/path", "example.com"))
}

func TestHasHost_SubdomainMatch(t *testing.T) {
	assert.True(t, HasHost("https://api.example.com/path", "example.com"))
}

func TestHasHost_PortStripped(t *testing.T) {
	assert.True(t, HasHost("https://example.com:8080/path", "example.com"))
}

func TestHasHost_MultipleHostnames(t *testing.T) {
	assert.True(t, HasHost("https://other.com/path", "example.com", "other.com"))
}

func TestHasHost_NoMatch(t *testing.T) {
	assert.False(t, HasHost("https://other.com/path", "example.com"))
}

func TestHasHost_InvalidURL(t *testing.T) {
	assert.False(t, HasHost("://invalid", "example.com"))
}

// endregion

// region - MergeMetadata

func TestMergeMetadata_OverlappingKeys(t *testing.T) {
	original, _ := types.NewMedia("https://example.com/a.jpg", types.Generic, map[string]interface{}{
		"name": "original",
		"id":   "123",
	}, nil)
	expanded, _ := types.NewMedia("https://example.com/b.jpg", types.Generic, map[string]interface{}{
		"name":   "expanded",
		"source": "test",
	}, nil)

	result := MergeMetadata(original, expanded)
	assert.Equal(t, "original", result.Metadata["name"])
	assert.Equal(t, "123", result.Metadata["id"])
	assert.Equal(t, "test", result.Metadata["source"])
}

func TestMergeMetadata_DisjointKeys(t *testing.T) {
	original, _ := types.NewMedia("https://example.com/a.jpg", types.Generic, map[string]interface{}{
		"name": "orig",
	}, nil)
	expanded, _ := types.NewMedia("https://example.com/b.jpg", types.Generic, map[string]interface{}{
		"source": "test",
	}, nil)

	result := MergeMetadata(original, expanded)
	assert.Equal(t, "orig", result.Metadata["name"])
	assert.Equal(t, "test", result.Metadata["source"])
}

func TestMergeMetadata_EmptyOriginal(t *testing.T) {
	original, _ := types.NewMedia("https://example.com/a.jpg", types.Generic, nil, nil)
	expanded, _ := types.NewMedia("https://example.com/b.jpg", types.Generic, map[string]interface{}{
		"source": "test",
	}, nil)

	result := MergeMetadata(original, expanded)
	assert.Equal(t, "test", result.Metadata["source"])
}

// endregion

// region - FilterMedia

func TestFilterMedia_PassThroughWhenEmpty(t *testing.T) {
	mediaCh := make(chan types.Media, 2)
	out := make(chan saktypes.Result[types.Media], 2)

	m1, _ := types.NewMedia("https://example.com/image.jpg", types.Generic, nil, nil)
	m2, _ := types.NewMedia("https://example.com/video.mp4", types.Generic, nil, nil)
	mediaCh <- m1
	mediaCh <- m2
	close(mediaCh)

	FilterMedia(mediaCh, nil, out)
	close(out)

	results := make([]types.Media, 0)
	for r := range out {
		results = append(results, r.Data)
	}
	assert.Equal(t, 2, len(results))
}

func TestFilterMedia_FilterByExtension(t *testing.T) {
	mediaCh := make(chan types.Media, 2)
	out := make(chan saktypes.Result[types.Media], 2)

	m1, _ := types.NewMedia("https://example.com/image.jpg", types.Generic, nil, nil)
	m2, _ := types.NewMedia("https://example.com/video.mp4", types.Generic, nil, nil)
	mediaCh <- m1
	mediaCh <- m2
	close(mediaCh)

	FilterMedia(mediaCh, []string{"jpg"}, out)
	close(out)

	results := make([]types.Media, 0)
	for r := range out {
		results = append(results, r.Data)
	}
	assert.Equal(t, 1, len(results))
	assert.Equal(t, "jpg", results[0].Extension)
}

func TestFilterMedia_NoMatches(t *testing.T) {
	mediaCh := make(chan types.Media, 2)
	out := make(chan saktypes.Result[types.Media], 2)

	m1, _ := types.NewMedia("https://example.com/image.jpg", types.Generic, nil, nil)
	mediaCh <- m1
	close(mediaCh)

	FilterMedia(mediaCh, []string{"png"}, out)
	close(out)

	results := make([]types.Media, 0)
	for r := range out {
		results = append(results, r.Data)
	}
	assert.Equal(t, 0, len(results))
}

// endregion

// region - FakeTimestamp

func TestFakeTimestamp_Deterministic(t *testing.T) {
	t1 := FakeTimestamp("test-input")
	t2 := FakeTimestamp("test-input")
	assert.Equal(t, t1, t2)
}

func TestFakeTimestamp_WithinRange(t *testing.T) {
	start := time.Date(1980, 10, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2035, 10, 1, 0, 0, 0, 0, time.UTC)

	ts := FakeTimestamp("some-string")
	assert.True(t, ts.After(start) || ts.Equal(start))
	assert.True(t, ts.Before(end) || ts.Equal(end))
}

func TestFakeTimestamp_DifferentInputsDiffer(t *testing.T) {
	t1 := FakeTimestamp("input-a")
	t2 := FakeTimestamp("input-b")
	assert.NotEqual(t, t1, t2)
}

// endregion
