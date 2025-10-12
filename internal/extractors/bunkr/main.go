package bunkr

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	saktypes "github.com/vegidio/go-sak/types"
	"github.com/vegidio/umd/internal/types"
	"github.com/vegidio/umd/internal/utils"
)

type Bunkr struct {
	Metadata types.Metadata

	url              string
	source           types.SourceType
	responseMetadata types.Metadata
	external         types.External
}

func New(url string, metadata types.Metadata, external types.External) (types.Extractor, error) {
	switch {
	case utils.HasHost(url, "bunkr.ac", "bunkr.ci", "bunkr.cr", "bunkr.fi", "bunkr.ph", "bunkr.pk",
		"bunkr.ps", "bunkr.si", "bunkr.sk", "bunkr.ws", "bunkr.black", "bunkr.red", "bunkr.media", "bunkr.site",
		"bunkr.ru"):
		return &Bunkr{Metadata: metadata, url: url, external: external}, nil
	}

	return nil, nil
}

func (b *Bunkr) Type() types.ExtractorType {
	return types.Bunkr
}

func (b *Bunkr) SourceType() (types.SourceType, error) {
	regexMedia1 := regexp.MustCompile(`/[fv]/([^/]+)/?$`)
	regexMedia2 := regexp.MustCompile(`cdn.+/([^/]+)/?$`)
	regexAlbum := regexp.MustCompile(`/a/([^/]+)/?$`)

	var source types.SourceType

	switch {
	case regexMedia1.MatchString(b.url):
		matches := regexMedia1.FindStringSubmatch(b.url)
		source = SourceMedia{id: matches[1]}
	case regexMedia2.MatchString(b.url):
		matches := regexMedia2.FindStringSubmatch(b.url)
		source = SourceMedia{id: matches[1]}
	case regexAlbum.MatchString(b.url):
		matches := regexAlbum.FindStringSubmatch(b.url)
		source = SourceAlbum{id: matches[1]}
	}

	if source == nil {
		return nil, fmt.Errorf("source type not found for URL: %s", b.url)
	}

	b.source = source
	return source, nil
}

func (b *Bunkr) QueryMedia(limit int, extensions []string, deep bool) (*types.Response, func()) {
	var err error
	ctx, stop := context.WithCancel(context.Background())

	if b.responseMetadata == nil {
		b.responseMetadata = make(types.Metadata)
	}

	response := &types.Response{
		Url:       b.url,
		Media:     make([]types.Media, 0),
		Extractor: types.Bunkr,
		Metadata:  b.responseMetadata,
		Done:      make(chan error),
	}

	go func() {
		defer close(response.Done)

		if b.source == nil {
			b.source, err = b.SourceType()
			if err != nil {
				response.Done <- err
				return
			}
		}

		mediaCh := b.fetchMedia(b.source, extensions, deep)

		for {
			select {
			case <-ctx.Done():
				return

			case result, ok := <-mediaCh:
				if !ok {
					return
				}

				if result.Err != nil {
					response.Done <- result.Err
					return
				}

				// Limiting the number of results
				if utils.MergeMedia(&response.Media, result.Data) >= limit {
					response.Media = response.Media[:limit]
					return
				}
			}
		}
	}()

	return response, stop
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

		out <- types.NewMedia(img.Url, types.Bunkr, map[string]interface{}{
			"id":      img.Slug,
			"name":    img.Name,
			"source":  strings.ToLower(sourceName),
			"created": img.Published,
		}, headers)
	}()

	return out
}

// endregion
