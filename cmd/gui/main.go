package main

import (
	"embed"
	"io"
	"shared"

	log "github.com/sirupsen/logrus"
	"github.com/vegidio/go-sak/o11y"
	"github.com/wailsapp/wails/v2/pkg/options/mac"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Disable logging
	log.SetOutput(io.Discard)

	otel := o11y.NewTelemetry(
		shared.OtelEndpoint,
		"umd",
		shared.Version,
		map[string]string{"Authorization": shared.OtelAuth},
		shared.OtelEnvironment,
		true,
	)

	defer otel.Close()

	// Create an instance of the app structure
	app := NewApp(otel)

	// Create an application with options
	err := wails.Run(&options.App{
		Title:  "Universal Media Downloader",
		Width:  960,
		Height: 710, // ideally 720
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		Mac: &mac.Options{
			WebviewIsTransparent: true,
		},
		BackgroundColour: &options.RGBA{R: 1, G: 1, B: 1, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		otel.LogError("Error running the app", nil, err)
		log.Error("Error:", err)
	}
}
