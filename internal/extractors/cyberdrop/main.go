package cyberdrop

import (
	"fmt"
	"regexp"
	"strings"

	saktypes "github.com/vegidio/go-sak/types"
	"github.com/vegidio/umd/internal/types"
	"github.com/vegidio/umd/internal/utils"
)

type Cyberdrop struct {
	types.BaseExtractor
}

func New(url string, metadata types.Metadata, external types.External) (types.Extractor, error) {
	switch {
	case utils.HasHost(url, "cyberdrop.to", "cyberdrop.me", "cyberdrop.cr"):
		c := &Cyberdrop{}
		c.BaseExtractor = types.BaseExtractor{
			Metadata: metadata,
			Url:      url,
			External: external,
			ExtType:  types.Cyberdrop,
		}
		c.FetchMediaFn = c.fetchMedia
		c.SourceTypeFn = c.SourceType
		return c, nil
	}

	return nil, nil
}

var (
	regexImage = regexp.MustCompile(`/f/([^/]+)/?$`)
	regexAlbum = regexp.MustCompile(`/a/([^/]+)/?$`)
)

func (c *Cyberdrop) SourceType() (types.SourceType, error) {

	var source types.SourceType

	switch {
	case regexImage.MatchString(c.Url):
		matches := regexImage.FindStringSubmatch(c.Url)
		source = SourceMedia{id: matches[1]}
	case regexAlbum.MatchString(c.Url):
		matches := regexAlbum.FindStringSubmatch(c.Url)
		source = SourceAlbum{id: matches[1]}
	}

	if source == nil {
		return nil, fmt.Errorf("source type not found for URL: %s", c.Url)
	}

	c.Source = source
	return source, nil
}

func (c *Cyberdrop) DownloadHeaders() map[string]string {
	return nil
}

// Compile-time assertion to ensure the extractor implements the Extractor interface
var _ types.Extractor = (*Cyberdrop)(nil)

// region - Private methods

func (c *Cyberdrop) fetchMedia(
	source types.SourceType,
	_ int,
	extensions []string,
	_ bool,
) <-chan saktypes.Result[types.Media] {
	out := make(chan saktypes.Result[types.Media])

	go func() {
		defer close(out)
		var images <-chan saktypes.Result[Image]

		switch s := source.(type) {
		case SourceMedia:
			images = c.fetchImage(s)
		case SourceAlbum:
			images = c.fetchAlbum(s)
		}

		for img := range images {
			if img.Err != nil {
				out <- saktypes.Result[types.Media]{Err: img.Err}
				return
			}

			media := c.dataToMedia(img.Data, source.Type())
			utils.FilterMedia(media, extensions, out)
		}
	}()

	return out
}

func (c *Cyberdrop) fetchImage(source SourceMedia) <-chan saktypes.Result[Image] {
	out := make(chan saktypes.Result[Image])

	go func() {
		defer close(out)
		img, err := getImage(source.id)

		if err != nil {
			out <- saktypes.Result[Image]{Err: err}
		} else {
			out <- saktypes.Result[Image]{Data: *img}
		}
	}()

	return out
}

func (c *Cyberdrop) fetchAlbum(source SourceAlbum) <-chan saktypes.Result[Image] {
	out := make(chan saktypes.Result[Image])

	go func() {
		defer close(out)

		ids, err := getAlbum(source.id)
		if err != nil {
			out <- saktypes.Result[Image]{Err: err}
			return
		}

		for _, id := range ids {
			img, iErr := getImage(id)

			if iErr != nil {
				out <- saktypes.Result[Image]{Err: iErr}
			} else {
				out <- saktypes.Result[Image]{Data: *img}
			}
		}
	}()

	return out
}

func (c *Cyberdrop) dataToMedia(img Image, sourceName string) <-chan types.Media {
	out := make(chan types.Media)
	headers := c.DownloadHeaders()

	go func() {
		defer close(out)

		tempUrl := "https://cyberdrop.to/" + img.Name
		media, err := types.NewMedia(tempUrl, types.Cyberdrop, map[string]interface{}{
			"id":      img.Id,
			"name":    img.Name,
			"source":  strings.ToLower(sourceName),
			"created": img.Published,
		}, headers)
		if err != nil {
			return
		}

		media.Url = img.Url
		out <- media
	}()

	return out
}

// endregion
