package shared

import (
	"sync"

	"github.com/samber/lo"
	"github.com/vegidio/go-sak/fetch"
)

type Queue struct {
	items []*fetch.Response
	mu    sync.RWMutex
}

func NewQueue(capacity int) *Queue {
	return &Queue{
		items: make([]*fetch.Response, 0, capacity),
	}
}

func (q *Queue) Items() []*fetch.Response {
	q.mu.RLock()
	defer q.mu.RUnlock()

	return q.items
}

func (q *Queue) Incompleted() int {
	q.mu.RLock()
	defer q.mu.RUnlock()

	return lo.CountBy(q.items, func(item *fetch.Response) bool {
		return !item.IsComplete()
	})
}

func (q *Queue) Add(item *fetch.Response) {
	q.mu.Lock()
	defer q.mu.Unlock()

	incomplete := lo.Filter(q.items, func(item *fetch.Response, _ int) bool {
		return !item.IsComplete()
	})

	q.items = append(incomplete, item)
}
