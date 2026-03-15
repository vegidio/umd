package simpcity

import (
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
	types.BaseExtractor
}

func New(url string, metadata types.Metadata, external types.External) (types.Extractor, error) {
	switch {
	case utils.HasHost(url, "simpcity.cr"):
		s := &SimpCity{}
		s.BaseExtractor = types.BaseExtractor{
			Metadata: metadata,
			Url:      url,
			External: external,
			ExtType:  types.SimpCity,
		}
		s.FetchMediaFn = s.fetchMedia
		s.SourceTypeFn = s.SourceType

		cookie, exists := metadata[types.SimpCity]["cookie"].(string)
		if !exists || len(cookie) == 0 {
			return s, fmt.Errorf("the extractor SimpCity requires cookies")
		}

		return s, nil
	}

	return nil, nil
}

var regexThread = regexp.MustCompile(`/threads/([^/]+)/?.*$`)

func (s *SimpCity) SourceType() (types.SourceType, error) {

	var source types.SourceType

	switch {
	case regexThread.MatchString(s.Url):
		matches := regexThread.FindStringSubmatch(s.Url)
		source = SourceThread{id: matches[1]}
	}

	if source == nil {
		return nil, fmt.Errorf("source type not found for URL: %s", s.Url)
	}

	s.Source = source
	return source, nil
}

func (s *SimpCity) DownloadHeaders() map[string]string {
	headers := make(map[string]string)

	cookie, exists := s.Metadata[types.SimpCity]["cookie"].(string)
	if exists {
		headers["Cookie"] = cookie
	}

	return headers
}

// Compile-time assertion to ensure the extractor implements the Extractor interface
var _ types.Extractor = (*SimpCity)(nil)

// region - Private methods

func (s *SimpCity) fetchMedia(
	source types.SourceType,
	_ int,
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

			media := s.dataToMedia(post.Data, source.Type())
			if deep {
				media = async.ConcurrentChannel(media, 5, func(m types.Media) types.Media {
					return s.External.ExpandMedia(m, Host, &s.ResponseMetadata)
				})
			}

			utils.FilterMedia(media, extensions, out)
		}
	}()

	return out
}

func (s *SimpCity) dataToMedia(post Post, sourceName string) <-chan types.Media {
	out := make(chan types.Media)
	headers := s.DownloadHeaders()

	go func() {
		defer close(out)

		for _, attachment := range post.Attachments {
			media, err := types.NewMedia(attachment.ThumbUrl, types.SimpCity, map[string]interface{}{
				"name":    post.Id,
				"url":     post.Url,
				"title":   post.Title,
				"source":  strings.ToLower(sourceName),
				"created": post.Published,
			}, headers)
			if err != nil {
				continue
			}

			media.Url = attachment.MediaUrl
			out <- media
		}
	}()

	return out
}

// endregion
