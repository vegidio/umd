package charm

import (
	"io"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"github.com/vegidio/go-sak/fetch"
)

// finishedResponse builds a Response that has already run to completion. The Progress is left at 0 to
// mimic a download that failed: those never reach 100%, which is exactly what used to keep the
// progress bar below 100% forever (see issue #24).
func finishedResponse(url string, filePath string, progress float64) *fetch.Response {
	done := make(chan struct{})
	close(done)

	return &fetch.Response{
		Request:  &fetch.Request{Url: url, FilePath: filePath},
		Progress: progress,
		Done:     done,
	}
}

// runProgress drives the real Bubble Tea loop and reports whether it terminated on its own.
func runProgress(t *testing.T, responses []*fetch.Response, total int) bool {
	t.Helper()

	ch := make(chan *fetch.Response, len(responses))
	for _, r := range responses {
		ch <- r
	}
	close(ch)

	program := tea.NewProgram(
		initProgressModel(ch, total),
		tea.WithInput(nil),
		tea.WithOutput(io.Discard),
	)

	exited := make(chan struct{})
	go func() {
		defer close(exited)
		_, _ = program.Run()
	}()

	select {
	case <-exited:
		return true
	case <-time.After(10 * time.Second):
		program.Kill()
		return false
	}
}

func TestProgress_ExitsWhenEveryDownloadSucceeds(t *testing.T) {
	responses := []*fetch.Response{
		finishedResponse("https://example.com/a.jpg", "/tmp/a.jpg", 1),
		finishedResponse("https://example.com/b.jpg", "/tmp/b.jpg", 1),
	}

	assert.True(t, runProgress(t, responses, len(responses)), "the progress UI should exit on its own")
}

// Regression test for https://github.com/vegidio/umd/issues/24: a download that ends in an error never
// reports 100% progress, and the UI used to spin forever at 99.x% because of it.
func TestProgress_ExitsWhenSomeDownloadsFail(t *testing.T) {
	responses := []*fetch.Response{
		finishedResponse("https://example.com/a.jpg", "/tmp/a.jpg", 1),
		finishedResponse("/r/some-subreddit", "/tmp/unknown-1.unknown", 0),
		finishedResponse("https://www.deviantart.com/someone", "/tmp/unknown-2.unknown", 0),
	}

	assert.True(t, runProgress(t, responses, len(responses)), "failed downloads must still count as finished")
}

func TestProgress_ExitsWhenThereIsNoMedia(t *testing.T) {
	assert.True(t, runProgress(t, nil, 0), "an empty media list should not hang the progress UI")
}
