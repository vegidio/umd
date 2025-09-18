package utils

import "github.com/vegidio/umd/internal/types"

func MergeMetadata(originalMedia types.Media, expandedMedia types.Media) types.Media {
	for k, v := range originalMedia.Metadata {
		expandedMedia.Metadata[k] = v
	}

	return expandedMedia
}
