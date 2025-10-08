package utils

import (
	"github.com/samber/lo"
	"github.com/vegidio/umd/internal/types"
)

func MergeMedia(media *[]types.Media, newMedia types.Media) int {
	*media = lo.UniqBy(append(*media, newMedia), func(m types.Media) string {
		return m.Url
	})

	return len(*media)
}
