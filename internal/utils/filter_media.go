package utils

import (
	"context"
	"slices"

	saktypes "github.com/vegidio/go-sak/types"
	"github.com/vegidio/umd/internal/types"
)

// FilterMedia filters media items based on the provided extensions.
// If the extensions' list is empty, all media items are returned.
func FilterMedia(
	ctx context.Context,
	media <-chan types.Media,
	extensions []string,
	out chan<- saktypes.Result[types.Media],
) {
	for m := range media {
		// Filter files with certain extensions; an empty list lets everything through
		if len(extensions) > 0 && !slices.Contains(extensions, m.Extension) {
			continue
		}

		if !Send(ctx, out, saktypes.Result[types.Media]{Data: m}) {
			return
		}
	}
}
