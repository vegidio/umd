package jpgfish

import (
	"fmt"
	"regexp"
	"strings"

	saktypes "github.com/vegidio/go-sak/types"
	"github.com/vegidio/umd/internal/types"
	"github.com/vegidio/umd/internal/utils"
)

type JpgFish struct {
	types.BaseExtractor
}

func New(url string, metadata types.Metadata, external types.External) (types.Extractor, error) {
	switch {
	case utils.HasHost(url, "jpg5.su", "jpg6.su", "jpg7.cr"):
		j := &JpgFish{}
		j.BaseExtractor = types.BaseExtractor{
			Metadata: metadata,
			Url:      url,
			External: external,
			ExtType:  types.JpgFish,
		}
		j.FetchMediaFn = j.fetchMedia
		j.SourceTypeFn = j.SourceType
		return j, nil
	}

	return nil, nil
}

var (
	regexImage = regexp.MustCompile(`/img/([^/]+)/?$`)
	regexAlbum = regexp.MustCompile(`/a/([^/]+)/?$`)
	regexUser  = regexp.MustCompile(`/([^/]+)/?$`)
)

func (j *JpgFish) SourceType() (types.SourceType, error) {

	var source types.SourceType

	switch {
	case regexImage.MatchString(j.Url):
		matches := regexImage.FindStringSubmatch(j.Url)
		source = SourceImage{id: matches[1]}
	case regexAlbum.MatchString(j.Url):
		matches := regexAlbum.FindStringSubmatch(j.Url)
		source = SourceAlbum{id: matches[1]}
	case regexUser.MatchString(j.Url):
		matches := regexUser.FindStringSubmatch(j.Url)
		source = SourceUser{name: matches[1]}
	}

	if source == nil {
		return nil, fmt.Errorf("source type not found for URL: %s", j.Url)
	}

	j.Source = source
	return source, nil
}

func (j *JpgFish) DownloadHeaders() map[string]string {
	return nil
}

// Compile-time assertion to ensure the extractor implements the Extractor interface
var _ types.Extractor = (*JpgFish)(nil)

// region - Private methods

func (j *JpgFish) fetchMedia(
	source types.SourceType,
	_ int,
	extensions []string,
	_ bool,
) <-chan saktypes.Result[types.Media] {
	out := make(chan saktypes.Result[types.Media])

	go func() {
		defer close(out)
		var images <-chan saktypes.Result[Image]

		switch s := source.(type) {
		case SourceImage:
			images = j.fetchImage(s)
		case SourceAlbum, SourceUser:
			return
		}

		for img := range images {
			if img.Err != nil {
				out <- saktypes.Result[types.Media]{Err: img.Err}
				return
			}

			media := j.dataToMedia(img.Data, source.Type())
			utils.FilterMedia(media, extensions, out)
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
		} else {
			result <- saktypes.Result[Image]{Data: *img}
		}
	}()

	return result
}

// endregion

// region - Private functions

func (j *JpgFish) dataToMedia(img Image, sourceName string) <-chan types.Media {
	out := make(chan types.Media)
	headers := j.DownloadHeaders()

	go func() {
		defer close(out)

		media, err := types.NewMedia(img.Url, types.JpgFish, map[string]interface{}{
			"id":      img.Id,
			"name":    img.Author,
			"title":   img.Title,
			"source":  strings.ToLower(sourceName),
			"created": img.Published,
		}, headers)
		if err != nil {
			return
		}
		out <- media
	}()

	return out
}

// endregion
