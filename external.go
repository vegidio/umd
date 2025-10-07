package umd

import (
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/vegidio/umd/internal/types"
	"github.com/vegidio/umd/internal/utils"
)

type External struct{}

func (External) ExpandMedia(media []Media, ignoreHost string, metadata *Metadata, parallel int) []Media {
	result := make([]Media, 0)

	var mu sync.Mutex
	var wg sync.WaitGroup
	sem := make(chan struct{}, parallel)

	for _, m := range media {
		wg.Add(1)

		go func(current Media) {
			defer func() {
				<-sem
				wg.Done()
			}()

			sem <- struct{}{}

			if current.Type == types.Unknown && !utils.HasHost(current.Url, ignoreHost) {
				extractor, err := New().
					WithMetadata(*metadata).
					FindExtractor(current.Url)

				if err != nil {
					appendResult(&mu, &result, current)
					return
				}

				log.WithFields(log.Fields{
					"url": current.Url,
				}).Debug("Expanding media")

				resp, _ := extractor.QueryMedia(1, nil, false)
				if resp.Error() != nil {
					appendResult(&mu, &result, current)
					return
				}

				mu.Lock()
				if _, exists := (*metadata)[resp.Extractor]; !exists {
					(*metadata)[resp.Extractor] = resp.Metadata[resp.Extractor]
				}
				mu.Unlock()

				if len(resp.Media) > 0 {
					mu.Lock()
					resp.Media[0] = utils.MergeMetadata(m, resp.Media[0])
					result = append(result, resp.Media[0])
					mu.Unlock()
				}
			} else {
				appendResult(&mu, &result, current)
			}
		}(m)
	}

	wg.Wait()
	close(sem)

	return result
}

func appendResult(mu *sync.Mutex, result *[]Media, media Media) {
	mu.Lock()
	*result = append(*result, media)
	mu.Unlock()
}
