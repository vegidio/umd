package simpcity

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

type SimpCity struct {
	Metadata types.Metadata

	url              string
	source           types.SourceType
	responseMetadata types.Metadata
	headers          map[string]string
	external         types.External
}

func New(url string, metadata types.Metadata, headers map[string]string, external types.External) types.Extractor {
	switch {
	case utils.HasHost(url, "simpcity.cr"):
		return &SimpCity{Metadata: metadata, url: url, headers: headers, external: external}
	}

	return nil
}

func (s *SimpCity) Type() types.ExtractorType {
	return types.SimpCity
}

func (s *SimpCity) SourceType() (types.SourceType, error) {
	regexThread := regexp.MustCompile(`/threads/([^/]+)/?$`)

	var source types.SourceType

	switch {
	case regexThread.MatchString(s.url):
		matches := regexThread.FindStringSubmatch(s.url)
		source = SourceThread{id: matches[1]}
	}

	if source == nil {
		return nil, fmt.Errorf("source type not found for URL: %s", s.url)
	}

	s.source = source
	return source, nil
}

func (s *SimpCity) QueryMedia(limit int, extensions []string, deep bool) (*types.Response, func()) {
	var err error
	ctx, stop := context.WithCancel(context.Background())

	if s.responseMetadata == nil {
		s.responseMetadata = make(types.Metadata)
	}

	response := &types.Response{
		Url:       s.url,
		Media:     make([]types.Media, 0),
		Extractor: types.SimpCity,
		Metadata:  s.responseMetadata,
		Done:      make(chan error),
	}

	go func() {
		defer close(response.Done)

		if s.source == nil {
			s.source, err = s.SourceType()
			if err != nil {
				response.Done <- err
				return
			}
		}

		mediaCh := s.fetchMedia(s.source, extensions, deep)

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
var _ types.Extractor = (*SimpCity)(nil)

// region - Private methods

func (s *SimpCity) fetchMedia(
	source types.SourceType,
	extensions []string,
	_ bool,
) <-chan saktypes.Result[[]types.Media] {
	out := make(chan saktypes.Result[[]types.Media])

	maxPages, exists := s.Metadata[types.SimpCity]["maxPages"].(int)
	if !exists {
		maxPages = 0
	}

	go func() {
		defer close(out)
		var posts <-chan saktypes.Result[Post]

		switch ss := source.(type) {
		case SourceThread:
			posts = getThread(ss.id, maxPages, s.headers)
		}

		for post := range posts {
			if post.Err != nil {
				out <- saktypes.Result[[]types.Media]{Err: post.Err}
				return
			}

			media := postToMedia(post.Data, source.Type())

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

// endregion

// region - Private functions

func postToMedia(post Post, sourceName string) []types.Media {
	media := make([]types.Media, 0)

	attachmentMedia := lo.Map(post.Attachments, func(attachment Attachment, index int) types.Media {
		newMedia := types.NewMedia(attachment.ThumbUrl, types.SimpCity, map[string]interface{}{
			"id":      post.Id,
			"url":     post.Url,
			"name":    post.Name,
			"source":  strings.ToLower(sourceName),
			"created": post.Published,
		})

		newMedia.Url = attachment.MediaUrl
		return newMedia
	})

	media = append(media, attachmentMedia...)
	return media
}

// endregion
