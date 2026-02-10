package services

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// MediaServer runs a local HTTP server dedicated to serving media files (video).
// This bypasses Wails' AssetServer / WebKitGTK custom URI scheme, which does not
// properly support video streaming (Range requests, proper Content-Type, etc.).
type MediaServer struct {
	mu       sync.Mutex
	listener net.Listener
	port     int
}

func NewMediaServer() *MediaServer {
	return &MediaServer{}
}

// Start begins listening on a random localhost port.  Call from OnStartup.
func (m *MediaServer) Start() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.listener != nil {
		return nil // already running
	}

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return fmt.Errorf("media server listen: %w", err)
	}

	m.listener = ln
	m.port = ln.Addr().(*net.TCPAddr).Port

	cwd, _ := os.Getwd()
	imagesDir := filepath.Join(cwd, "images")

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Only allow GET
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Sanitise: only serve from imagesDir, no directory traversal
		name := filepath.Base(r.URL.Path)
		if name == "." || name == "/" || strings.Contains(name, "..") {
			http.NotFound(w, r)
			return
		}

		target := filepath.Join(imagesDir, name)
		// Ensure target is inside imagesDir
		if rel, err := filepath.Rel(imagesDir, target); err != nil || strings.HasPrefix(rel, "..") {
			http.NotFound(w, r)
			return
		}

		// http.ServeFile on a real TCP listener supports Range / 206 correctly.
		http.ServeFile(w, r, target)
	})

	go func() {
		if err := http.Serve(ln, mux); err != nil && !isClosedErr(err) {
			log.Printf("[MediaServer] serve error: %v", err)
		}
	}()

	log.Printf("[MediaServer] listening on 127.0.0.1:%d", m.port)
	return nil
}

// Stop shuts down the listener.  Call from OnShutdown.
func (m *MediaServer) Stop(_ context.Context) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.listener != nil {
		m.listener.Close()
		m.listener = nil
	}
}

// Port returns the listening port.  Exposed to frontend via Wails binding.
func (m *MediaServer) Port() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.port
}

func isClosedErr(err error) bool {
	return strings.Contains(err.Error(), "use of closed network connection")
}
