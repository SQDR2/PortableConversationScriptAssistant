package main

import (
	"context"
	"fmt"
	"os"
	"strings"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// GetVersion returns the application version from the VERSION file
func (a *App) GetVersion() string {
	version, err := os.ReadFile("VERSION")
	if err != nil {
		return "v1.0.0" // 默认回退版本
	}
	return strings.TrimSpace(string(version))
}
