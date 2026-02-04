package main

import (
	"context"
	"embed"
	"log"
	"os"
	"sidekick/backend/db"
	"sidekick/backend/services"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	if os.Getenv("XDG_SESSION_TYPE") == "wayland" || os.Getenv("WAYLAND_DISPLAY") != "" {
		// Force X11 backend for reliable window positioning on Linux
		_ = os.Setenv("GDK_BACKEND", "x11")
		log.Printf("Detected Wayland session; forcing GDK_BACKEND=x11 for window positioning")
	}

	// Create an instance of the app structure
	app := NewApp()
	windowService := services.NewWindowService()
	scriptService := services.NewScriptService()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "sidekick",
		Width:  350,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup: func(ctx context.Context) {
			app.startup(ctx)
			windowService.Startup(ctx)
			scriptService.Startup(ctx)

			// Initialize DB
			// TODO: Get proper app data dir
			cwd, _ := os.Getwd()
			if err := db.InitDB(cwd); err != nil {
				runtime.LogErrorf(ctx, "Failed to init DB: %v", err)
			}
		},
		OnShutdown: func(ctx context.Context) {
			windowService.Shutdown(ctx)
		},
		Bind: []interface{}{
			app,
			windowService,
			scriptService,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
