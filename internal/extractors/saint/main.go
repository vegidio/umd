package saint

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	saktypes "github.com/vegidio/go-sak/types"
	"github.com/vegidio/umd/internal/types"
	"github.com/vegidio/umd/internal/utils"
)

type Saint struct {
	Metadata types.Metadata

	url              string
	source           types.SourceType
	responseMetadata types.Metadata
	external         types.External
}

func New(url string, metadata types.Metadata, external types.External) (types.Extractor, error) {
	switch {
	case utils.HasHost(url, "saint.to", "saint2.su", "saint2.cr"):
		return &Saint{Metadata: metadata, url: url, external: external}, nil
	}

	return nil, nil
}

func (s *Saint) Type() types.ExtractorType {
	return types.Saint
}

func (s *Saint) SourceType() (types.SourceType, error) {
	regexVideo := regexp.MustCompile(`/embed/([^/]+)/?$`)

	var source types.SourceType

	switch {
	case regexVideo.MatchString(s.url):
		matches := regexVideo.FindStringSubmatch(s.url)
		source = SourceVideo{id: matches[1]}
	}

	if source == nil {
		return nil, fmt.Errorf("source type not found for URL: %s", s.url)
	}

	s.source = source
	return source, nil
}

func (s *Saint) QueryMedia(limit int, extensions []string, deep bool) (*types.Response, func()) {
	var err error
	ctx, stop := context.WithCancel(context.Background())

	if s.responseMetadata == nil {
		s.responseMetadata = make(types.Metadata)
	}

	response := &types.Response{
		Url:       s.url,
		Media:     make([]types.Media, 0),
		Extractor: types.Saint,
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

func (s *Saint) DownloadHeaders() map[string]string {
	return nil
}

// Compile-time assertion to ensure the extractor implements the Extractor interface
var _ types.Extractor = (*Saint)(nil)

// region - Private methods

func (s *Saint) fetchMedia(
	source types.SourceType,
	extensions []string,
	_ bool,
) <-chan saktypes.Result[types.Media] {
	out := make(chan saktypes.Result[types.Media])

	go func() {
		defer close(out)
		var videos <-chan saktypes.Result[Video]

		switch ss := source.(type) {
		case SourceVideo:
			videos = s.fetchVideo(ss)
		}

		for video := range videos {
			if video.Err != nil {
				out <- saktypes.Result[types.Media]{Err: video.Err}
				return
			}

			media := s.dataToMedia(video.Data, source.Type())
			utils.FilterMedia(media, extensions, out)
		}
	}()

	return out
}

func (s *Saint) fetchVideo(source SourceVideo) <-chan saktypes.Result[Video] {
	result := make(chan saktypes.Result[Video])

	go func() {
		defer close(result)
		video, err := getVideo(source.id)

		if err != nil {
			result <- saktypes.Result[Video]{Err: err}
		} else {
			result <- saktypes.Result[Video]{Data: *video}
		}
	}()

	return result
}

func (s *Saint) dataToMedia(video Video, sourceName string) <-chan types.Media {
	out := make(chan types.Media)
	headers := s.DownloadHeaders()

	go func() {
		defer close(out)

		out <- types.NewMedia(video.Url, types.Saint, map[string]interface{}{
			"id":      video.Id,
			"name":    video.Id,
			"source":  strings.ToLower(sourceName),
			"created": video.Published,
		}, headers)
	}()

	return out
}

// endregion
