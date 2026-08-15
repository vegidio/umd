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
	types.BaseExtractor
}

func New(url string, metadata types.Metadata, external types.External) (types.Extractor, error) {
	switch {
	case utils.HasHost(url, "saint.to", "saint2.su", "saint2.cr", "turbo.cr"):
		s := &Saint{}
		s.BaseExtractor = types.BaseExtractor{
			Metadata: metadata,
			Url:      url,
			External: external,
			ExtType:  types.Saint,
		}
		s.FetchMediaFn = s.fetchMedia
		s.SourceTypeFn = s.SourceType
		return s, nil
	}

	return nil, nil
}

var regexVideo = regexp.MustCompile(`/embed/([^/]+)/?$`)

func (s *Saint) SourceType() (types.SourceType, error) {
	var source types.SourceType

	switch {
	case regexVideo.MatchString(s.Url):
		matches := regexVideo.FindStringSubmatch(s.Url)
		source = SourceVideo{id: matches[1]}
	}

	if source == nil {
		return nil, fmt.Errorf("source type not found for URL: %s", s.Url)
	}

	s.Source = source
	return source, nil
}

func (s *Saint) DownloadHeaders() map[string]string {
	return nil
}

// Compile-time assertion to ensure the extractor implements the Extractor interface
var _ types.Extractor = (*Saint)(nil)

// region - Private methods

func (s *Saint) fetchMedia(
	ctx context.Context,
	source types.SourceType,
	_ int,
	extensions []string,
	_ bool,
) <-chan saktypes.Result[types.Media] {
	out := make(chan saktypes.Result[types.Media])

	go func() {
		defer close(out)
		var videos <-chan saktypes.Result[Video]

		switch ss := source.(type) {
		case SourceVideo:
			videos = s.fetchVideo(ctx, ss)
		}

		for video := range videos {
			if video.Err != nil {
				utils.Send(ctx, out, saktypes.Result[types.Media]{Err: video.Err})
				return
			}

			media := s.dataToMedia(ctx, video.Data, source.Type())
			utils.FilterMedia(ctx, media, extensions, out)
		}
	}()

	return out
}

func (s *Saint) fetchVideo(ctx context.Context, source SourceVideo) <-chan saktypes.Result[Video] {
	result := make(chan saktypes.Result[Video])

	go func() {
		defer close(result)
		video, err := getVideo(ctx, source.id)

		if err != nil {
			utils.Send(ctx, result, saktypes.Result[Video]{Err: err})
			return
		}

		utils.Send(ctx, result, saktypes.Result[Video]{Data: *video})
	}()

	return result
}

func (s *Saint) dataToMedia(ctx context.Context, video Video, sourceName string) <-chan types.Media {
	out := make(chan types.Media)
	headers := s.DownloadHeaders()

	go func() {
		defer close(out)

		media, err := types.NewMedia(video.Url, types.Saint, map[string]interface{}{
			"id":      video.Id,
			"name":    video.Id,
			"source":  strings.ToLower(sourceName),
			"created": video.Published,
		}, headers)
		if err != nil {
			return
		}
		utils.Send(ctx, out, media)
	}()

	return out
}

// endregion
