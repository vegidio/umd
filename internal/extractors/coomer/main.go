package coomer

import (
	"context"
	"fmt"
	"reflect"
	"regexp"
	"slices"

	"github.com/samber/lo"
	saktypes "github.com/vegidio/go-sak/types"
	"github.com/vegidio/umd/internal/types"
	"github.com/vegidio/umd/internal/utils"
)

type Coomer struct {
	Metadata types.Metadata

	url              string
	extractor        types.ExtractorType
	source           types.SourceType
	services         string
	responseMetadata types.Metadata
	headers          map[string]string
	external         types.External
}

func New(url string, metadata types.Metadata, headers map[string]string, external types.External) types.Extractor {
	switch {
	case utils.HasHost(url, "coomer.party", "coomer.st", "coomer.su"):
		baseUrl = "https://coomer.st"

		return &Coomer{
			Metadata: metadata,

			url:       url,
			extractor: types.Coomer,
			services:  "onlyfans|fansly|candfans",
			headers:   headers,
			external:  external,
		}
	case utils.HasHost(url, "kemono.party", "kemono.su", "kemono.cr"):
		baseUrl = "https://kemono.cr"

		return &Coomer{
			Metadata: metadata,

			url:       url,
			extractor: types.Kemono,
			services:  "patreon|fanbox|discord|fantia|afdian|boosty|gumroad|subscribestar|dlsite",
			headers:   headers,
			external:  external,
		}
	}

	return nil
}

func (c *Coomer) Type() types.ExtractorType {
	return c.extractor
}

func (c *Coomer) SourceType() (types.SourceType, error) {
	regexPost := regexp.MustCompile(`(` + c.services + `)/user/([^/]+)/post/([^/\n?]+)`)
	regexUser := regexp.MustCompile(`(` + c.services + `)/user/([^/\n?]+)`)

	var source types.SourceType
	var user string

	switch {
	case regexPost.MatchString(c.url):
		matches := regexPost.FindStringSubmatch(c.url)
		service := matches[1]
		user = matches[2]
		id := matches[3]
		source = SourcePost{Service: service, Id: id, name: user}

	case regexUser.MatchString(c.url):
		matches := regexUser.FindStringSubmatch(c.url)
		service := matches[1]
		user = matches[2]
		source = SourceUser{Service: service, name: user}
	}

	if source == nil {
		return nil, fmt.Errorf("source type not found for URL: %s", c.url)
	}

	c.source = source
	return source, nil
}

func (c *Coomer) QueryMedia(limit int, extensions []string, deep bool) (*types.Response, func()) {
	var err error
	ctx, stop := context.WithCancel(context.Background())

	if c.responseMetadata == nil {
		c.responseMetadata = make(types.Metadata)
	}

	response := &types.Response{
		Url:       c.url,
		Media:     make([]types.Media, 0),
		Extractor: c.extractor,
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

// Compile-time assertion to ensure the extractor implements the Extractor interface
var _ types.Extractor = (*Coomer)(nil)

// region - Private methods

func (c *Coomer) fetchMedia(
	source types.SourceType,
	extensions []string,
	_ bool,
) <-chan saktypes.Result[[]types.Media] {
	out := make(chan saktypes.Result[[]types.Media])

	sourceValue := reflect.ValueOf(source)
	serviceField := sourceValue.FieldByName("Service")
	profile, pErr := getProfile(serviceField.String(), source.Name())

	if pErr != nil {
		out <- saktypes.Result[[]types.Media]{Err: pErr}
		return out
	}

	go func() {
		defer close(out)
		var responses <-chan saktypes.Result[Response]

		switch s := source.(type) {
		case SourceUser:
			responses = getUser(*profile)
		case SourcePost:
			responses = getPost(*profile, s.Id)
		}

		for response := range responses {
			if response.Err != nil {
				out <- saktypes.Result[[]types.Media]{Err: response.Err}
				return
			}

			media := c.postToMedia(response.Data, profile.Name)

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

func (c *Coomer) postToMedia(response Response, name string) []types.Media {
	media := make([]types.Media, 0)

	for _, image := range response.Images {
		if image.Path != "" {
			url := image.Server + "/data" + image.Path
			newMedia := types.NewMedia(url, c.extractor, map[string]interface{}{
				"source":  response.Post.Service,
				"name":    name,
				"created": response.Post.Published.Time,
			})

			media = append(media, newMedia)
		}
	}

	for _, video := range response.Videos {
		if video.Path != "" {
			url := video.Server + "/data" + video.Path
			newMedia := types.NewMedia(url, c.extractor, map[string]interface{}{
				"source":  response.Post.Service,
				"name":    name,
				"created": response.Post.Published.Time,
			})

			media = append(media, newMedia)
		}
	}

	return media
}

// endregion
