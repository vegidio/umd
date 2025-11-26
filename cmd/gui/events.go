package main

import (
	"github.com/vegidio/go-sak/fetch"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func (a *App) OnExtractorFound(name string) {
	runtime.EventsEmit(a.ctx, "OnExtractorFound", name)
}

func (a *App) OnExtractorTypeFound(eType string, name string) {
	runtime.EventsEmit(a.ctx, "OnExtractorTypeFound", eType, name)
}

func (a *App) OnMediaQueried(amount int) {
	runtime.EventsEmit(a.ctx, "OnMediaQueried", amount)
}

func (a *App) OnQueryCompleted(total int, isCached bool) {
	runtime.EventsEmit(a.ctx, "OnQueryCompleted", total, isCached)
}

func (a *App) OnMediaDownloaded(amount int, responses []*fetch.Response) {
	runtime.EventsEmit(a.ctx, "OnMediaDownloaded", amount, responses)
}
