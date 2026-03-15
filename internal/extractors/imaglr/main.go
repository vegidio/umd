package imaglr

import (
	"fmt"
	"regexp"
	"strings"

	saktypes "github.com/vegidio/go-sak/types"
	"github.com/vegidio/umd/internal/types"
	"github.com/vegidio/umd/internal/utils"
)

type Imaglr struct {
	types.BaseExtractor
}

func New(url string, metadata types.Metadata, external types.External) (types.Extractor, error) {
	switch {
	case utils.HasHost(url, "imaglr.com"):
		i := &Imaglr{}
		i.BaseExtractor = types.BaseExtractor{
			Metadata: metadata,
			Url:      url,
			External: external,
			ExtType:  types.Imaglr,
		}
		i.FetchMediaFn = i.fetchMedia
		i.SourceTypeFn = i.SourceType
		return i, nil
	}

	return nil, nil
}

var regexPost = regexp.MustCompile(`/post/([^/\n?]+)`)

func (i *Imaglr) SourceType() (types.SourceType, error) {

	var source types.SourceType
	var id string

	switch {
	case regexPost.MatchString(i.Url):
		matches := regexPost.FindStringSubmatch(i.Url)
		id = matches[1]
		source = SourcePost{name: id}
	}

	if source == nil {
		return nil, fmt.Errorf("source type not found for URL: %s", i.Url)
	}

	i.Source = source
	return source, nil
}

func (i *Imaglr) DownloadHeaders() map[string]string {
	return nil
}

// Compile-time assertion to ensure the extractor implements the Extractor interface
var _ types.Extractor = (*Imaglr)(nil)

// region - Private methods

func (i *Imaglr) fetchMedia(
	source types.SourceType,
	_ int,
	extensions []string,
	_ bool,
) <-chan saktypes.Result[types.Media] {
	out := make(chan saktypes.Result[types.Media])

	go func() {
		defer close(out)

		posts := make([]Post, 0)
		var err error

		switch s := source.(type) {
		case SourcePost:
			posts, err = i.fetchPost(s)
		}

		if err != nil {
			out <- saktypes.Result[types.Media]{Err: err}
			return
		}

		media := i.dataToMedia(posts, source.Name())
		utils.FilterMedia(media, extensions, out)
	}()

	return out
}

func (i *Imaglr) fetchPost(source SourcePost) ([]Post, error) {
	post, err := getPost(source.name)

	if err != nil {
		return make([]Post, 0), err
	}

	return []Post{*post}, nil
}

func (i *Imaglr) dataToMedia(posts []Post, sourceName string) <-chan types.Media {
	out := make(chan types.Media)
	headers := i.DownloadHeaders()

	go func() {
		defer close(out)

		for _, post := range posts {
			media, err := types.NewMedia(post.Media, types.Imaglr, map[string]interface{}{
				"id":      post.Id,
				"name":    post.Author,
				"source":  strings.ToLower(sourceName),
				"created": post.Timestamp,
			}, headers)
			if err != nil {
				continue
			}
			out <- media
		}
	}()

	return out
}

// endregion
