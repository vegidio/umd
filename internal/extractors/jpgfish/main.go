package jpgfish

import (
	"context"
	"fmt"
	"regexp"
	"slices"
	"strings"

	"github.com/samber/lo"
	saktypes "github.com/vegidio/go-sak/types"
	"github.com/vegidio/umd/internal/types"
	"github.com/vegidio/umd/internal/utils"
)

type JpgFish struct {
	Metadata types.Metadata

	url              string
	source           types.SourceType
	responseMetadata types.Metadata
	external         types.External
}

func New(url string, metadata types.Metadata, external types.External) types.Extractor {
	switch {
	case utils.HasHost(url, "jpg5.su") || utils.HasHost(url, "jpg6.su"):
		return &JpgFish{Metadata: metadata, url: url, external: external}
	}

	return nil
}

func (j *JpgFish) Type() types.ExtractorType {
	return types.JpgFish
}

func (j *JpgFish) SourceType() (types.SourceType, error) {
	regexImage := regexp.MustCompile(`/img/([^/]+)/?$`)
	regexUser := regexp.MustCompile(`/([^/]+)/?$`)

	var source types.SourceType

	switch {
	case regexImage.MatchString(j.url):
		matches := regexImage.FindStringSubmatch(j.url)
		source = SourceImage{id: matches[1]}
	case regexUser.MatchString(j.url):
		matches := regexImage.FindStringSubmatch(j.url)
		source = SourceUser{name: matches[1]}
	}

	if source == nil {
		return nil, fmt.Errorf("source type not found for URL: %s", j.url)
	}

	j.source = source
	return source, nil
}

func (j *JpgFish) QueryMedia(limit int, extensions []string, deep bool) (*types.Response, func()) {
	var err error
	ctx, stop := context.WithCancel(context.Background())

	if j.responseMetadata == nil {
		j.responseMetadata = make(types.Metadata)
	}

	response := &types.Response{
		Url:       j.url,
		Media:     make([]types.Media, 0),
		Extractor: types.JpgFish,
		Metadata:  j.responseMetadata,
		Done:      make(chan error),
	}

	go func() {
		defer close(response.Done)

		if j.source == nil {
			j.source, err = j.SourceType()
			if err != nil {
				response.Done <- err
				return
			}
		}

		mediaCh := j.fetchMedia(j.source, limit, extensions, deep)

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

// Compile-time assertion to ensure the extractor implements the Extractor interface
var _ types.Extractor = (*JpgFish)(nil)

// region - Private methods

func (j *JpgFish) fetchMedia(
	source types.SourceType,
	limit int,
	extensions []string,
	_ bool,
) <-chan saktypes.Result[[]types.Media] {
	out := make(chan saktypes.Result[[]types.Media])

	go func() {
		defer close(out)
		var images <-chan saktypes.Result[Image]

		switch s := source.(type) {
		case SourceImage:
			images = j.fetchImage(s)
		}

		for img := range images {
			if img.Err != nil {
				out <- saktypes.Result[[]types.Media]{Err: img.Err}
				return
			}

			media := imageToMedia(img.Data, source.Type())

			// Filter files with certain extensions
			if len(extensions) > 0 {
				media = lo.Filter(media, func(m types.Media, _ int) bool {
					return slices.Contains(extensions, m.Extension)
				})
			}

			out <- saktypes.Result[[]types.Media]{Data: media}
		}
	}()

	return out
}

func (j *JpgFish) fetchImage(source SourceImage) <-chan saktypes.Result[Image] {
	result := make(chan saktypes.Result[Image])

	go func() {
		defer close(result)

		img, err := getImage(source.id)
		if err != nil {
			result <- saktypes.Result[Image]{Err: err}
		}

		result <- saktypes.Result[Image]{Data: *img}
	}()

	return result
}

// endregion

// region - Private functions

func imageToMedia(img Image, sourceName string) []types.Media {
	return []types.Media{types.NewMedia(img.Url, types.JpgFish, map[string]interface{}{
		"id":      img.Id,
		"name":    img.Author,
		"title":   img.Title,
		"source":  strings.ToLower(sourceName),
		"created": img.Published,
	})}
}

// endregion
