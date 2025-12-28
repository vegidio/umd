package main

import (
	"context"

	"github.com/vegidio/go-sak/o11y"
)

// App struct
type App struct {
	ctx context.Context
	tel *o11y.Telemetry
}

// NewApp creates a new App application struct
func NewApp(tel *o11y.Telemetry) *App {
	return &App{tel: tel}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}
