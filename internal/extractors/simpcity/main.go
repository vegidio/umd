package simpcity

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/vegidio/go-sak/async"
	saktypes "github.com/vegidio/go-sak/types"
	"github.com/vegidio/umd/internal/types"
	"github.com/vegidio/umd/internal/utils"
)

const Host = "simpcity.cr"

type SimpCity struct {
	Metadata types.Metadata

	url              string
	source           types.SourceType
	responseMetadata types.Metadata
	external         types.External
}

func New(url string, metadata types.Metadata, external types.External) (types.Extractor, error) {
	switch {
	case utils.HasHost(url, "simpcity.cr"):
		ext := &SimpCity{Metadata: metadata, url: url, external: external}

		cookie, exists := metadata[types.SimpCity]["cookie"].(string)
		if !exists || len(cookie) == 0 {
			return ext, fmt.Errorf("the extractor SimpCity requires cookies")
		}

		return ext, nil
	}

	return nil, nil
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
	deep bool,
) <-chan saktypes.Result[types.Media] {
	out := make(chan saktypes.Result[types.Media])

	headers := make(map[string]string)
	cookie, exists := s.Metadata[types.SimpCity]["cookie"].(string)
	if exists {
		headers["Cookie"] = cookie
	}

	startPage, exists := s.Metadata[types.SimpCity]["startPage"].(int)
	if !exists {
		startPage = 1
	}

	maxPages, exists := s.Metadata[types.SimpCity]["maxPages"].(int)
	if !exists {
		maxPages = 0
	}

	go func() {
		defer close(out)
		var posts <-chan saktypes.Result[Post]

		switch ss := source.(type) {
		case SourceThread:
			posts = getThread(ss.id, startPage, maxPages, headers)
		}

		for post := range posts {
			if post.Err != nil {
				out <- saktypes.Result[types.Media]{Err: post.Err}
				return
			}

			media := postToMedia(post.Data, source.Type())
			if deep {
				media = async.ProcessChannel(media, 5, func(m types.Media) types.Media {
					return s.external.ExpandMedia(m, Host, &s.responseMetadata)
				})
			}

			utils.FilterMedia(media, extensions, out)
		}
	}()

	return out
}

// endregion

// region - Private functions

func postToMedia(post Post, sourceName string) <-chan types.Media {
	out := make(chan types.Media)

	go func() {
		defer close(out)

		for _, attachment := range post.Attachments {
			media := types.NewMedia(attachment.ThumbUrl, types.SimpCity, map[string]interface{}{
				"id":      post.Id,
				"url":     post.Url,
				"name":    post.Name,
				"source":  strings.ToLower(sourceName),
				"created": post.Published,
			})

			media.Url = attachment.MediaUrl
			out <- media
		}
	}()

	return out
}

// endregion
