package types

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	saktypes "github.com/vegidio/go-sak/types"
)

type fakeSource struct{}

func (fakeSource) Type() string { return "fake" }
func (fakeSource) Name() string { return "fake" }

func newFakeExtractor(fetch FetchMediaFunc) *BaseExtractor {
	return &BaseExtractor{
		Url:          "https://example.com",
		ExtType:      Generic,
		SourceTypeFn: func() (SourceType, error) { return fakeSource{}, nil },
		FetchMediaFn: fetch,
	}
}

// The cancel function used to stop only the consumer, leaving every request in flight and every
// producer goroutine blocked on its next send.
func TestQueryMediaContext_StopCancelsInFlightRequest(t *testing.T) {
	cancelled := make(chan struct{}, 1)

	server := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		select {
		case <-r.Context().Done():
			cancelled <- struct{}{}
		case <-time.After(10 * time.Second):
		}
	}))
	defer server.Close()

	extractor := newFakeExtractor(func(
		ctx context.Context,
		_ SourceType,
		_ int,
		_ []string,
		_ bool,
	) <-chan saktypes.Result[Media] {
		out := make(chan saktypes.Result[Media])

		go func() {
			defer close(out)

			req, _ := http.NewRequestWithContext(ctx, http.MethodGet, server.URL, nil)
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				return
			}

			_ = resp.Body.Close()
		}()

		return out
	})

	response, stop := extractor.QueryMediaContext(context.Background(), 10, nil, false)

	time.Sleep(200 * time.Millisecond)
	stop()

	select {
	case <-cancelled:
	case <-time.After(5 * time.Second):
		t.Fatal("stop() did not cancel the in-flight request")
	}

	queryDone := make(chan struct{})
	go func() {
		defer close(queryDone)
		_ = response.Error()
	}()

	select {
	case <-queryDone:
	case <-time.After(5 * time.Second):
		t.Fatal("the query did not finish after stop()")
	}
}

// Reaching the limit ends the query, and the producer must be released rather than left blocked on a
// channel nobody reads anymore.
func TestQueryMediaContext_ReachingTheLimitReleasesTheProducer(t *testing.T) {
	producerDone := make(chan struct{})

	extractor := newFakeExtractor(func(
		ctx context.Context,
		_ SourceType,
		_ int,
		_ []string,
		_ bool,
	) <-chan saktypes.Result[Media] {
		out := make(chan saktypes.Result[Media])

		go func() {
			defer close(out)
			defer close(producerDone)

			for i := 0; ; i++ {
				media, err := NewMedia(fmt.Sprintf("https://example.com/%d.jpg", i), Generic, nil, nil)
				if err != nil {
					return
				}

				select {
				case out <- saktypes.Result[Media]{Data: media}:
				case <-ctx.Done():
					return
				}
			}
		}()

		return out
	})

	response, _ := extractor.QueryMediaContext(context.Background(), 3, nil, false)
	require.NoError(t, response.Error())
	assert.Len(t, response.Media, 3)

	select {
	case <-producerDone:
	case <-time.After(5 * time.Second):
		t.Fatal("the producer was left blocked after the limit was reached")
	}
}
