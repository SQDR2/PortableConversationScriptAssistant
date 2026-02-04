//go:build linux

package utils

import (
	"fmt"
	"log"
	"strings"

	"github.com/jezek/xgb/xproto"
	"github.com/jezek/xgbutil"
	"github.com/jezek/xgbutil/ewmh"
	"github.com/jezek/xgbutil/icccm"
	"github.com/jezek/xgbutil/xwindow"
)

type LinuxProvider struct {
	X *xgbutil.XUtil
}

func NewWindowProvider() WindowProvider {
	X, err := xgbutil.NewConn()
	if err != nil {
		log.Printf("Failed to connect to X: %v", err)
		return &LinuxProvider{}
	}
	return &LinuxProvider{X: X}
}

func (p *LinuxProvider) Close() error {
	if p.X != nil {
		p.X.Conn().Close()
	}
	return nil
}

func (p *LinuxProvider) GetWindows() ([]WindowInfo, error) {
	if p.X == nil {
		return nil, fmt.Errorf("X connection not available")
	}

	clientIds, err := ewmh.ClientListGet(p.X)
	if err != nil {
		return nil, err
	}

	var windows []WindowInfo
	for _, id := range clientIds {
		name, err := ewmh.WmNameGet(p.X, id)
		if err != nil || name == "" {
			name, _ = icccm.WmNameGet(p.X, id)
		}

		if name == "" {
			continue // Skip unnamed windows
		}

		states, _ := ewmh.WmStateGet(p.X, id)
		isIconic := false
		for _, s := range states {
			if s == "_NET_WM_STATE_HIDDEN" {
				isIconic = true
				break
			}
		}

		geom, err := xwindow.New(p.X, id).Geometry()
		if err != nil {
			continue
		}

		// Use WM_CLASS for process identification context
		wmClass, _ := icccm.WmClassGet(p.X, id)
		processName := ""
		if wmClass != nil {
			processName = wmClass.Instance + " (" + wmClass.Class + ")"
		}

		windows = append(windows, WindowInfo{
			Handle:  fmt.Sprintf("%d", id),
			Title:   name,
			Process: processName,
			Rect: Rect{
				Left:   geom.X(),
				Top:    geom.Y(),
				Right:  geom.X() + geom.Width(),
				Bottom: geom.Y() + geom.Height(),
			},
			IsIconic: isIconic,
		})
	}
	return windows, nil
}

func (p *LinuxProvider) GetWindowRect(handleStr string) (Rect, bool, error) {
	if p.X == nil {
		return Rect{}, false, fmt.Errorf("X connection not available")
	}

	var id uint32
	fmt.Sscanf(handleStr, "%d", &id)
	wid := xproto.Window(id)

	geom, err := xwindow.New(p.X, wid).Geometry()
	if err != nil {
		return Rect{}, false, err
	}

	states, _ := ewmh.WmStateGet(p.X, wid)
	isIconic := false
	for _, s := range states {
		if s == "_NET_WM_STATE_HIDDEN" {
			isIconic = true
			break
		}
	}

	return Rect{
		Left:   geom.X(),
		Top:    geom.Y(),
		Right:  geom.X() + geom.Width(),
		Bottom: geom.Y() + geom.Height(),
	}, isIconic, nil
}

func ForceMoveResizeWindowByTitle(title string, x int, y int, width int, height int) error {
	X, err := xgbutil.NewConn()
	if err != nil {
		return fmt.Errorf("failed to connect to X: %w", err)
	}
	defer X.Conn().Close()

	clientIds, err := ewmh.ClientListGet(X)
	if err != nil {
		return err
	}

	for _, id := range clientIds {
		name, err := ewmh.WmNameGet(X, id)
		if err != nil || name == "" {
			name, _ = icccm.WmNameGet(X, id)
		}
		if !strings.Contains(name, title) {
			continue
		}
		win := xwindow.New(X, id)
		if width > 0 && height > 0 {
			win.Resize(width, height)
		}
		win.Move(x, y)
		return nil
	}

	return fmt.Errorf("window not found: %s", title)
}

func GetWindowDecorationHeightByTitle(title string) (int, error) {
	X, err := xgbutil.NewConn()
	if err != nil {
		return 0, fmt.Errorf("failed to connect to X: %w", err)
	}
	defer X.Conn().Close()

	clientIds, err := ewmh.ClientListGet(X)
	if err != nil {
		return 0, err
	}

	for _, id := range clientIds {
		name, err := ewmh.WmNameGet(X, id)
		if err != nil || name == "" {
			name, _ = icccm.WmNameGet(X, id)
		}
		if !strings.Contains(name, title) {
			continue
		}

		_ = ewmh.RequestFrameExtents(X, id)
		extents, err := ewmh.FrameExtentsGet(X, id)
		if err != nil || extents == nil {
			return 0, fmt.Errorf("failed to get frame extents: %s", name)
		}
		decorationHeight := int(extents.Top + extents.Bottom)
		if decorationHeight < 0 {
			decorationHeight = 0
		}
		return decorationHeight, nil
	}

	return 0, fmt.Errorf("window not found: %s", title)
}
