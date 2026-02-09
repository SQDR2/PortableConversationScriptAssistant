package main

import (
	"context"
	"embed"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sidekick/backend/db"
	"sidekick/backend/services"
	"sidekick/backend/utils"
	"strings"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Initialize crash logging (must be first)
	utils.InitLogDir()
	defer func() {
		if r := recover(); r != nil {
			utils.LogError("panic", r)
		}
	}()

	if os.Getenv("XDG_SESSION_TYPE") == "wayland" || os.Getenv("WAYLAND_DISPLAY") != "" {
		// Force X11 backend for reliable window positioning on Linux
		_ = os.Setenv("GDK_BACKEND", "x11")
		log.Printf("Detected Wayland session; forcing GDK_BACKEND=x11 for window positioning")
	}

	// Create an instance of the app structure
	app := NewApp()
	windowService := services.NewWindowService()
	scriptService := services.NewScriptService()
	categoryService := services.NewCategoryService()
	translationService := services.NewTranslationService()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "sidekick",
		Width:  350,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if strings.HasPrefix(r.URL.Path, "/images/") {
					cwd, _ := os.Getwd()
					filePath := filepath.Join(cwd, r.URL.Path)
					if _, err := os.Stat(filePath); err == nil {
						http.ServeFile(w, r, filePath)
						return
					}
				}
				http.NotFound(w, r)
			}),
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup: func(ctx context.Context) {
			// Initialize DB
			// TODO: Get proper app data dir
			cwd, _ := os.Getwd()
			if err := db.InitDB(cwd); err != nil {
				runtime.LogErrorf(ctx, "Failed to init DB: %v", err)
				utils.LogError("startup_db", err)
			}

			app.startup(ctx)
			windowService.Startup(ctx)
			scriptService.Startup(ctx)
			categoryService.Startup(ctx)
			translationService.Startup(ctx)
		},
		OnShutdown: func(ctx context.Context) {
			windowService.Shutdown(ctx)
		},
		Bind: []interface{}{
			app,
			windowService,
			scriptService,
			categoryService,
			translationService,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
