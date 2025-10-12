package cyberdrop

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	saktypes "github.com/vegidio/go-sak/types"
	"github.com/vegidio/umd/internal/types"
	"github.com/vegidio/umd/internal/utils"
)

type Cyberdrop struct {
	Metadata types.Metadata

	url              string
	source           types.SourceType
	responseMetadata types.Metadata
	external         types.External
}

func New(url string, metadata types.Metadata, external types.External) (types.Extractor, error) {
	switch {
	case utils.HasHost(url, "cyberdrop.to", "cyberdrop.me"):
		return &Cyberdrop{Metadata: metadata, url: url, external: external}, nil
	}

	return nil, nil
}

func (c *Cyberdrop) Type() types.ExtractorType {
	return types.Cyberdrop
}

func (c *Cyberdrop) SourceType() (types.SourceType, error) {
	regexImage := regexp.MustCompile(`/f/([^/]+)/?$`)
	regexAlbum := regexp.MustCompile(`/a/([^/]+)/?$`)

	var source types.SourceType

	switch {
	case regexImage.MatchString(c.url):
		matches := regexImage.FindStringSubmatch(c.url)
		source = SourceMedia{id: matches[1]}
	case regexAlbum.MatchString(c.url):
		matches := regexAlbum.FindStringSubmatch(c.url)
		source = SourceAlbum{id: matches[1]}
	}

	if source == nil {
		return nil, fmt.Errorf("source type not found for URL: %s", c.url)
	}

	c.source = source
	return source, nil
}

func (c *Cyberdrop) QueryMedia(limit int, extensions []string, deep bool) (*types.Response, func()) {
	var err error
	ctx, stop := context.WithCancel(context.Background())

	if c.responseMetadata == nil {
		c.responseMetadata = make(types.Metadata)
	}

	response := &types.Response{
		Url:       c.url,
		Media:     make([]types.Media, 0),
		Extractor: types.Cyberdrop,
		Metadata:  c.responseMetadata,
		Done:      make(chan error),
	}

	go func() {
		defer close(response.Done)

		if c.source == nil {
			c.source, err = c.SourceType()
			if err != nil {
				response.Done <- err
				return
			}
		}

		mediaCh := c.fetchMedia(c.source, extensions, deep)

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

func (c *Cyberdrop) DownloadHeaders() map[string]string {
	return nil
}

// Compile-time assertion to ensure the extractor implements the Extractor interface
var _ types.Extractor = (*Cyberdrop)(nil)

// region - Private methods

func (c *Cyberdrop) fetchMedia(
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
		media := types.NewMedia(tempUrl, types.Cyberdrop, map[string]interface{}{
			"id":      img.Id,
			"name":    img.Name,
			"source":  strings.ToLower(sourceName),
			"created": img.Published,
		}, headers)

		media.Url = img.Url
		out <- media
	}()

	return out
}

// endregion
