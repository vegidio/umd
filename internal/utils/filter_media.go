package utils

import (
	"slices"

	saktypes "github.com/vegidio/go-sak/types"
	"github.com/vegidio/umd/internal/types"
)

// FilterMedia filters media items based on the provided extensions.
// If the extensions' list is empty, all media items are returned.
func FilterMedia(
	media <-chan types.Media,
	extensions []string,
	out chan<- saktypes.Result[types.Media],
) {
	for m := range media {
		// If no extensions specified, pass through all media
		if len(extensions) == 0 {
			out <- saktypes.Result[types.Media]{Data: m}
			continue
		}

		// Filter files with certain extensions
		if slices.Contains(extensions, m.Extension) {
			out <- saktypes.Result[types.Media]{Data: m}
		}
	}
}
