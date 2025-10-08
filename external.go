package umd

import (
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/vegidio/umd/internal/types"
	"github.com/vegidio/umd/internal/utils"
)

type External struct{}

func (External) ExpandMedia(media Media, ignoreHost string, metadata *Metadata) Media {
	var mu sync.Mutex

	if media.Type == types.Unknown && !utils.HasHost(media.Url, ignoreHost) {
		extractor, err := New().
			WithMetadata(*metadata).
			FindExtractor(media.Url)

		if err != nil {
			return media
		}

		log.WithFields(log.Fields{
			"url": media.Url,
		}).Debug("Expanding media")

		resp, _ := extractor.QueryMedia(1, nil, false)
		if resp.Error() != nil {
			return media
		}

		mu.Lock()
		if _, exists := (*metadata)[resp.Extractor]; !exists {
			(*metadata)[resp.Extractor] = resp.Metadata[resp.Extractor]
		}
		mu.Unlock()

		if len(resp.Media) > 0 {
			resp.Media[0] = utils.MergeMetadata(media, resp.Media[0])
			return resp.Media[0]
		}
	}

	return media
}
