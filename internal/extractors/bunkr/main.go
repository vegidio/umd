package bunkr

import (
	"fmt"
	"regexp"
	"strings"

	saktypes "github.com/vegidio/go-sak/types"
	"github.com/vegidio/umd/internal/types"
	"github.com/vegidio/umd/internal/utils"
)

type Bunkr struct {
	types.BaseExtractor
}

func New(url string, metadata types.Metadata, external types.External) (types.Extractor, error) {
	switch {
	case utils.HasHost(url, "bunkr.ac", "bunkr.ci", "bunkr.cr", "bunkr.fi", "bunkr.ph", "bunkr.pk",
		"bunkr.ps", "bunkr.si", "bunkr.sk", "bunkr.ws", "bunkr.black", "bunkr.red", "bunkr.media", "bunkr.site",
		"bunkr.ru"):
		b := &Bunkr{}
		b.BaseExtractor = types.BaseExtractor{
			Metadata: metadata,
			Url:      url,
			External: external,
			ExtType:  types.Bunkr,
		}
		b.FetchMediaFn = b.fetchMedia
		b.SourceTypeFn = b.SourceType
		return b, nil
	}

	return nil, nil
}

var (
	regexMedia1 = regexp.MustCompile(`/[fv]/([^/]+)/?$`)
	regexMedia2 = regexp.MustCompile(`cdn.+/([^/]+)/?$`)
	regexAlbum  = regexp.MustCompile(`/a/([^/]+)/?$`)
)

func (b *Bunkr) SourceType() (types.SourceType, error) {
	var source types.SourceType

	switch {
	case regexMedia1.MatchString(b.Url):
		matches := regexMedia1.FindStringSubmatch(b.Url)
		source = SourceMedia{id: matches[1]}
	case regexMedia2.MatchString(b.Url):
		matches := regexMedia2.FindStringSubmatch(b.Url)
		source = SourceMedia{id: matches[1]}
	case regexAlbum.MatchString(b.Url):
		matches := regexAlbum.FindStringSubmatch(b.Url)
		source = SourceAlbum{id: matches[1]}
	}

	if source == nil {
		return nil, fmt.Errorf("source type not found for URL: %s", b.Url)
	}

	b.Source = source
	return source, nil
}

func (b *Bunkr) DownloadHeaders() map[string]string {
	return map[string]string{
		"Referer": "https://bunkr.cr/",
	}
}

// Compile-time assertion to ensure the extractor implements the Extractor interface
var _ types.Extractor = (*Bunkr)(nil)

// region - Private methods

func (b *Bunkr) fetchMedia(
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
			images = b.fetchImage(s)
		case SourceAlbum:
			images = b.fetchAlbum(s)
		}

		for img := range images {
			if img.Err != nil {
				out <- saktypes.Result[types.Media]{Err: img.Err}
				return
			}

			media := b.dataToMedia(img.Data, source.Type())
			utils.FilterMedia(media, extensions, out)
		}
	}()

	return out
}

func (b *Bunkr) fetchImage(source SourceMedia) <-chan saktypes.Result[Image] {
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

func (b *Bunkr) fetchAlbum(source SourceAlbum) <-chan saktypes.Result[Image] {
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

func (b *Bunkr) dataToMedia(img Image, sourceName string) <-chan types.Media {
	out := make(chan types.Media)
	headers := b.DownloadHeaders()

	go func() {
		defer close(out)

		media, err := types.NewMedia(img.Url, types.Bunkr, map[string]interface{}{
			"id":      img.Slug,
			"name":    img.Name,
			"source":  strings.ToLower(sourceName),
			"created": img.Published,
		}, headers)
		if err != nil {
			return
		}
		out <- media
	}()

	return out
}

// endregion
