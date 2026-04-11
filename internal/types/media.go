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

	// Headers contains the headers required to download the media.
	Headers map[string]string
}

func (m Media) String() string {
	return fmt.Sprintf("{Url: %s, Extension: %s, Type: %s, Extractor: %s, Metadata: %v}",
		m.Url, m.Extension, m.Type, m.Extractor, m.Metadata)
}

func NewMedia(urlStr string, extractor ExtractorType, metadata map[string]interface{}, headers map[string]string) (Media, error) {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return Media{}, fmt.Errorf("error parsing URL %q: %w", urlStr, err)
	}

	parsedURL.RawQuery = ""
	cleanUrl := parsedURL.String()

	var extension string
	if IsBunkrRedirectURL(cleanUrl) {
		extension = ""
	} else {
		extension = getExtension(cleanUrl)
	}

	if metadata == nil {
		metadata = make(map[string]interface{})
	}

	return Media{
		Url:       urlStr,
		Extension: extension,
		Type:      GetType(extension),
		Extractor: extractor,
		Metadata:  metadata,
		Headers:   headers,
	}, nil
}

// region - Private functions

func getExtension(urStr string) string {
	u, err := url.Parse(urStr)
	if err != nil {
		return ""
	}

	ext := path.Ext(u.Path)
	if ext == "" {
		return ""
	}

	return strings.ToLower(ext[1:])
}

// IsBunkrRedirectURL returns true if the URL is a Bunkr page URL (not a direct media link).
// These URLs look like media files but are actually pages that redirect to the real media.
func IsBunkrRedirectURL(urlStr string) bool {
	return strings.Contains(urlStr, "bunkr") &&
		(strings.Contains(urlStr, "//cdn") ||
			strings.Contains(urlStr, "/f/") ||
			strings.Contains(urlStr, "/v/"))
}

func GetType(extension string) MediaType {
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
