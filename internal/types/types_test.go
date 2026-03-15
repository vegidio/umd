package types

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// region - NewMedia

func TestNewMedia_ValidURLWithExtension(t *testing.T) {
	m, _ := NewMedia("https://example.com/image.jpg", Generic, nil, nil)
	assert.Equal(t, "https://example.com/image.jpg", m.Url)
	assert.Equal(t, "jpg", m.Extension)
	assert.Equal(t, Image, m.Type)
	assert.Equal(t, Generic, m.Extractor)
}

func TestNewMedia_URLWithoutExtension(t *testing.T) {
	m, _ := NewMedia("https://example.com/somepath", Generic, nil, nil)
	assert.Equal(t, "", m.Extension)
	assert.Equal(t, Unknown, m.Type)
}

func TestNewMedia_URLWithQueryParams(t *testing.T) {
	m, _ := NewMedia("https://example.com/image.png?width=100&height=200", Generic, nil, nil)
	assert.Equal(t, "png", m.Extension)
	assert.Equal(t, Image, m.Type)
	// Original URL with query params is preserved
	assert.Equal(t, "https://example.com/image.png?width=100&height=200", m.Url)
}

func TestNewMedia_NilMetadataDefaultsToEmptyMap(t *testing.T) {
	m, _ := NewMedia("https://example.com/image.jpg", Generic, nil, nil)
	assert.NotNil(t, m.Metadata)
	assert.Equal(t, 0, len(m.Metadata))
}

func TestNewMedia_MetadataPreserved(t *testing.T) {
	meta := map[string]interface{}{"name": "test", "id": "123"}
	m, _ := NewMedia("https://example.com/image.jpg", Generic, meta, nil)
	assert.Equal(t, "test", m.Metadata["name"])
	assert.Equal(t, "123", m.Metadata["id"])
}

func TestNewMedia_HeadersPreserved(t *testing.T) {
	headers := map[string]string{"Referer": "https://example.com"}
	m, _ := NewMedia("https://example.com/image.jpg", Generic, nil, headers)
	assert.Equal(t, "https://example.com", m.Headers["Referer"])
}

// endregion

// region - getExtension (tested via NewMedia)

func TestGetExtension_ImageExtensions(t *testing.T) {
	tests := []struct {
		url string
		ext string
	}{
		{"https://example.com/file.jpg", "jpg"},
		{"https://example.com/file.jpeg", "jpeg"},
		{"https://example.com/file.png", "png"},
		{"https://example.com/file.gif", "gif"},
		{"https://example.com/file.webp", "webp"},
		{"https://example.com/file.avif", "avif"},
	}

	for _, tt := range tests {
		t.Run(tt.ext, func(t *testing.T) {
			m, _ := NewMedia(tt.url, Generic, nil, nil)
			assert.Equal(t, tt.ext, m.Extension)
		})
	}
}

func TestGetExtension_VideoExtensions(t *testing.T) {
	tests := []struct {
		url string
		ext string
	}{
		{"https://example.com/file.mp4", "mp4"},
		{"https://example.com/file.webm", "webm"},
		{"https://example.com/file.mkv", "mkv"},
		{"https://example.com/file.mov", "mov"},
		{"https://example.com/file.m4v", "m4v"},
	}

	for _, tt := range tests {
		t.Run(tt.ext, func(t *testing.T) {
			m, _ := NewMedia(tt.url, Generic, nil, nil)
			assert.Equal(t, tt.ext, m.Extension)
		})
	}
}

func TestGetExtension_NoExtension(t *testing.T) {
	m, _ := NewMedia("https://example.com/path/to/resource", Generic, nil, nil)
	assert.Equal(t, "", m.Extension)
}

func TestGetExtension_BunkrSpecialCase(t *testing.T) {
	// Bunkr CDN URLs should return empty extension
	m, _ := NewMedia("https://cdn.bunkr.cr/somefile.jpg", Generic, nil, nil)
	assert.Equal(t, "", m.Extension)

	// Bunkr /f/ URLs should return empty extension
	m2, _ := NewMedia("https://bunkr.cr/f/somefile.jpg", Generic, nil, nil)
	assert.Equal(t, "", m2.Extension)

	// Bunkr /v/ URLs should return empty extension
	m3, _ := NewMedia("https://bunkr.cr/v/somefile.mp4", Generic, nil, nil)
	assert.Equal(t, "", m3.Extension)
}

func TestGetExtension_UpperCaseNormalized(t *testing.T) {
	m, _ := NewMedia("https://example.com/file.JPG", Generic, nil, nil)
	assert.Equal(t, "jpg", m.Extension)
}

// endregion

// region - getType (tested via NewMedia)

func TestGetType_ImageType(t *testing.T) {
	m, _ := NewMedia("https://example.com/file.jpg", Generic, nil, nil)
	assert.Equal(t, Image, m.Type)
}

func TestGetType_VideoType(t *testing.T) {
	m, _ := NewMedia("https://example.com/file.mp4", Generic, nil, nil)
	assert.Equal(t, Video, m.Type)
}

func TestGetType_UnknownType(t *testing.T) {
	m, _ := NewMedia("https://example.com/file.pdf", Generic, nil, nil)
	assert.Equal(t, Unknown, m.Type)
}

func TestGetType_NoExtensionIsUnknown(t *testing.T) {
	m, _ := NewMedia("https://example.com/resource", Generic, nil, nil)
	assert.Equal(t, Unknown, m.Type)
}

// endregion

// region - ExtractorType.String()

func TestExtractorType_String(t *testing.T) {
	tests := []struct {
		et       ExtractorType
		expected string
	}{
		{Generic, "Generic"},
		{Bunkr, "Bunkr"},
		{Coomer, "Coomer"},
		{Cyberdrop, "Cyberdrop"},
		{Erome, "Erome"},
		{Fapello, "Fapello"},
		{Imaglr, "Imaglr"},
		{JpgFish, "JpgFish"},
		{Kemono, "Kemono"},
		{Reddit, "Reddit"},
		{RedGifs, "RedGifs"},
		{Saint, "Saint"},
		{SimpCity, "SimpCity"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.et.String())
		})
	}
}

func TestExtractorType_UnknownValue(t *testing.T) {
	et := ExtractorType(999)
	assert.Equal(t, "Unknown", et.String())
}

// endregion

// region - MediaType.String()

func TestMediaType_String(t *testing.T) {
	assert.Equal(t, "Image", Image.String())
	assert.Equal(t, "Video", Video.String())
	assert.Equal(t, "Unknown", Unknown.String())
}

// endregion

// region - Response.Error()

func TestResponse_Error_ReturnsNilOnSuccess(t *testing.T) {
	resp := &Response{Done: make(chan error, 1)}
	go func() {
		resp.Done <- nil
	}()
	assert.NoError(t, resp.Error())
}

func TestResponse_Error_ReturnsErrorOnFailure(t *testing.T) {
	resp := &Response{Done: make(chan error, 1)}
	go func() {
		resp.Done <- fmt.Errorf("something went wrong")
	}()
	err := resp.Error()
	assert.Error(t, err)
	assert.Equal(t, "something went wrong", err.Error())
}

// endregion

// region - Media.String()

func TestMedia_String(t *testing.T) {
	m, _ := NewMedia("https://example.com/image.jpg", Generic, nil, nil)
	s := m.String()
	assert.Contains(t, s, "https://example.com/image.jpg")
	assert.Contains(t, s, "jpg")
	assert.Contains(t, s, "Generic")
}

// endregion
