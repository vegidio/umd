package types

import (
	"fmt"
	"net/url"
	"path"
	"strings"
)

// Media represents a media object.
type Media struct {
	// Url is the URL of the media.
	Url string

	// Extension is the extension of the media file, derived from the URL.
	Extension string

	// Type is the type of media, determined based on the file extension.
	Type MediaType

	// Extractor is the extractor used to fetch the media.
	Extractor ExtractorType

	// Metadata contains metadata about the media. Default is an empty map.
	Metadata map[string]interface{}
}

func (m Media) String() string {
	return fmt.Sprintf("{Url: %s, Extension: %s, Type: %s, Extractor: %s, Metadata: %v}",
		m.Url, m.Extension, m.Type, m.Extractor, m.Metadata)
}

func NewMedia(urlStr string, extractor ExtractorType, metadata map[string]interface{}) Media {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		panic("Error parsing URL: " + err.Error())
	}

	parsedURL.RawQuery = ""
	cleanUrl := parsedURL.String()
	extension := getExtension(cleanUrl)

	if metadata == nil {
		metadata = make(map[string]interface{})
	}

	return Media{
		Url:       urlStr,
		Extension: extension,
		Type:      getType(extension),
		Extractor: extractor,
		Metadata:  metadata,
	}
}

// region - Private functions

func getExtension(urStr string) string {
	u, err := url.Parse(urStr)
	if err != nil {
		panic(err)
	}

	ext := path.Ext(u.Path)
	if ext == "" {
		return ""
	}

	return strings.ToLower(ext[1:])
}

func getType(extension string) MediaType {
	lowerExt := strings.ToLower(extension)

	switch lowerExt {
	case "avif", "gif", "jpg", "jpeg", "png", "webp":
		return Image
	case "gifv", "m4v", "mkv", "mov", "mp4", "webm":
		return Video
	default:
		return Unknown
	}
}

// endregion
