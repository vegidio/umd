package shared

import (
	"testing"
	"time"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/vegidio/umd"
)

// region - IsValidURL

func TestIsValidURL_ValidHTTPS(t *testing.T) {
	assert.True(t, IsValidURL("https://example.com"))
}

func TestIsValidURL_ValidHTTPWithPath(t *testing.T) {
	assert.True(t, IsValidURL("http://example.com/path?q=1"))
}

func TestIsValidURL_MissingScheme(t *testing.T) {
	assert.False(t, IsValidURL("example.com"))
}

func TestIsValidURL_MissingHost(t *testing.T) {
	assert.False(t, IsValidURL("https://"))
}

func TestIsValidURL_EmptyString(t *testing.T) {
	assert.False(t, IsValidURL(""))
}

// endregion

// region - GetHost

func TestGetHost_WithPort(t *testing.T) {
	assert.Equal(t, "example.com", GetHost("https://example.com:8080/path"))
}

func TestGetHost_WithoutPort(t *testing.T) {
	assert.Equal(t, "example.com", GetHost("https://example.com/path"))
}

func TestGetHost_Subdomain(t *testing.T) {
	assert.Equal(t, "api.example.com", GetHost("https://api.example.com/path"))
}

func TestGetHost_InvalidURL(t *testing.T) {
	assert.Equal(t, "", GetHost("://invalid"))
}

// endregion

// region - CreateTimestamp

func TestCreateTimestamp_KnownValue(t *testing.T) {
	// 1000000 in base36 = "lfls"
	result := CreateTimestamp(1000000)
	assert.Equal(t, "00lfls", result)
}

func TestCreateTimestamp_Zero(t *testing.T) {
	result := CreateTimestamp(0)
	assert.Equal(t, "000000", result)
}

func TestCreateTimestamp_Negative(t *testing.T) {
	result := CreateTimestamp(-1)
	assert.Equal(t, "0000-1", result)
}

// endregion

// region - CreateHashSuffix

func TestCreateHashSuffix_Deterministic(t *testing.T) {
	h1 := CreateHashSuffix("test-input")
	h2 := CreateHashSuffix("test-input")
	assert.Equal(t, h1, h2)
}

func TestCreateHashSuffix_FourChars(t *testing.T) {
	h := CreateHashSuffix("anything")
	assert.Equal(t, 4, len(h))
}

func TestCreateHashSuffix_DifferentInputsDiffer(t *testing.T) {
	h1 := CreateHashSuffix("input-a")
	h2 := CreateHashSuffix("input-b")
	assert.NotEqual(t, h1, h2)
}

// endregion

// region - GetMediaType

func TestGetMediaType_ImageExtensions(t *testing.T) {
	tests := []string{"file.jpg", "file.jpeg", "file.png", "file.gif", "file.webp", "file.avif"}
	for _, f := range tests {
		t.Run(f, func(t *testing.T) {
			assert.Equal(t, "image", GetMediaType(f))
		})
	}
}

func TestGetMediaType_VideoExtensions(t *testing.T) {
	tests := []string{"file.mp4", "file.webm", "file.mkv", "file.mov", "file.m4v"}
	for _, f := range tests {
		t.Run(f, func(t *testing.T) {
			assert.Equal(t, "video", GetMediaType(f))
		})
	}
}

func TestGetMediaType_UnknownExtension(t *testing.T) {
	assert.Equal(t, "unkwn", GetMediaType("file.pdf"))
}

// endregion

// region - CreateFilePath

func TestCreateFilePath_Format(t *testing.T) {
	media := umd.Media{
		Url:       "https://example.com/image.jpg",
		Extension: "jpg",
		Metadata: map[string]interface{}{
			"name":    "testname",
			"created": time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC),
		},
	}

	path := CreateFilePath("/downloads", media)
	assert.Contains(t, path, "/downloads/testname-")
	assert.Contains(t, path, ".jpg")
}

func TestCreateFilePath_EmptyExtension(t *testing.T) {
	media := umd.Media{
		Url:       "https://example.com/resource",
		Extension: "",
		Metadata: map[string]interface{}{
			"name":    "test",
			"created": time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC),
		},
	}

	path := CreateFilePath("/downloads", media)
	assert.Contains(t, path, ".unknown")
}

func TestCreateFilePath_StringTimestamp(t *testing.T) {
	media := umd.Media{
		Url:       "https://example.com/image.jpg",
		Extension: "jpg",
		Metadata: map[string]interface{}{
			"name":    "test",
			"created": "2024-01-15T12:00:00Z",
		},
	}

	path := CreateFilePath("/downloads", media)
	assert.Contains(t, path, "/downloads/test-")
	assert.Contains(t, path, ".jpg")
}

// endregion

// region - RemoveDuplicates

func TestRemoveDuplicates_NoDuplicates(t *testing.T) {
	// Use in-memory filesystem
	oldFs := fs
	fs = afero.NewMemMapFs()
	defer func() { fs = oldFs }()

	downloads := []Download{
		{Url: "https://a.com/1.jpg", FilePath: "/tmp/1.jpg", Hash: "aaa", IsSuccess: true},
		{Url: "https://a.com/2.jpg", FilePath: "/tmp/2.jpg", Hash: "bbb", IsSuccess: true},
	}

	deleted, remaining := RemoveDuplicates(downloads, nil)
	assert.Equal(t, 0, deleted)
	assert.Equal(t, 2, len(remaining))
}

func TestRemoveDuplicates_WithDuplicates(t *testing.T) {
	// Use in-memory filesystem
	oldFs := fs
	fs = afero.NewMemMapFs()
	defer func() { fs = oldFs }()

	// Create the files so Remove can work
	afero.WriteFile(fs, "/tmp/1.jpg", []byte("data"), 0644)
	afero.WriteFile(fs, "/tmp/2.jpg", []byte("data"), 0644)

	downloads := []Download{
		{Url: "https://a.com/1.jpg", FilePath: "/tmp/1.jpg", Hash: "same", IsSuccess: true},
		{Url: "https://a.com/2.jpg", FilePath: "/tmp/2.jpg", Hash: "same", IsSuccess: true},
	}

	var callbackCalled int
	deleted, remaining := RemoveDuplicates(downloads, func(d Download) {
		callbackCalled++
	})

	assert.Equal(t, 1, deleted)
	assert.Equal(t, 1, len(remaining))
	assert.Equal(t, 1, callbackCalled)
}

func TestRemoveDuplicates_EmptyHashFiltered(t *testing.T) {
	oldFs := fs
	fs = afero.NewMemMapFs()
	defer func() { fs = oldFs }()

	downloads := []Download{
		{Url: "https://a.com/1.jpg", FilePath: "/tmp/1.jpg", Hash: "", IsSuccess: true},
		{Url: "https://a.com/2.jpg", FilePath: "/tmp/2.jpg", Hash: "aaa", IsSuccess: true},
	}

	deleted, remaining := RemoveDuplicates(downloads, nil)
	assert.Equal(t, 0, deleted)
	// Only the one with a non-empty hash appears in remaining
	assert.Equal(t, 1, len(remaining))
}

// endregion
