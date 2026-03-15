package utils

import "github.com/vegidio/umd/internal/types"

func MergeMedia(media *[]types.Media, seen map[string]struct{}, newMedia types.Media) int {
	if _, exists := seen[newMedia.Url]; !exists {
		seen[newMedia.Url] = struct{}{}
		*media = append(*media, newMedia)
	}

	return len(*media)
}
