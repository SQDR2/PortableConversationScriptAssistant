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

func (p *LinuxProvider) GetForegroundHandle() string {
	if p.X == nil {
		return ""
	}
	id, err := ewmh.ActiveWindowGet(p.X)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%d", id)
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

func GetDWMFrameOffsetsByTitle(title string) (FrameOffsets, error) {
	// On Linux, X11 Geometry() and Resize() operate on the client area
	// (excluding WM decorations like title bar and borders).
	// We need to compensate: if sidekick has WM decorations, setting
	// client height = target height will make the visual window taller
	// than the target by the decoration amount.
	// Return negative Bottom offset = -(top + bottom extents) so that
	// the alignment formula: newHeight = targetHeight + Top + Bottom
	// becomes: newHeight = targetHeight - totalDecorationHeight.
	X, err := xgbutil.NewConn()
	if err != nil {
		return FrameOffsets{}, nil
	}
	defer X.Conn().Close()

	clientIds, err := ewmh.ClientListGet(X)
	if err != nil {
		return FrameOffsets{}, nil
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
			return FrameOffsets{}, nil
		}

		// Total decoration height that makes the visual window taller
		// than the client area set by Resize()
		totalDecor := int(extents.Top + extents.Bottom)
		return FrameOffsets{
			Top:    0,
			Bottom: -totalDecor, // negative: shrink client area to compensate
			Left:   0,
			Right:  0,
		}, nil
	}

	return FrameOffsets{}, nil
}

func GetWindowPhysicalWidthByTitle(title string) (int, error) {
	X, err := xgbutil.NewConn()
	if err != nil {
		return 0, nil
	}
	defer X.Conn().Close()

	clientIds, err := ewmh.ClientListGet(X)
	if err != nil {
		return 0, nil
	}

	for _, id := range clientIds {
		name, err := ewmh.WmNameGet(X, id)
		if err != nil || name == "" {
			name, _ = icccm.WmNameGet(X, id)
		}
		if !strings.Contains(name, title) {
			continue
		}

		geom, err := xwindow.New(X, id).Geometry()
		if err != nil {
			return 0, nil
		}

		// On Linux, Geometry().Width() is the client area width.
		// Add frame extents (left + right borders) to get visual width.
		totalWidth := geom.Width()
		_ = ewmh.RequestFrameExtents(X, id)
		extents, err := ewmh.FrameExtentsGet(X, id)
		if err == nil && extents != nil {
			totalWidth += int(extents.Left + extents.Right)
		}
		return totalWidth, nil
	}

	return 0, nil
}

func ForceRaiseWindowByTitle(title string) error {
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
		// 1. Physically raise the window in X server stack
		win.Stack(xproto.StackModeAbove)

		// 2. Add 'Above' state to satisfy EWMH window managers (like GNOME/KDE)
		_ = ewmh.WmStateReq(X, id, ewmh.StateAdd, "_NET_WM_STATE_ABOVE")

		return nil
	}

	return fmt.Errorf("window not found: %s", title)
}

func ForceLowerWindowByTitle(title string) error {
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

		// Remove 'Above' state
		_ = ewmh.WmStateReq(X, id, ewmh.StateRemove, "_NET_WM_STATE_ABOVE")

		// Gently push it below the current layer to ensure non-AlwaysOnTop windows (like App A)
		// can cover it once focus is lost.
		win := xwindow.New(X, id)
		win.Stack(xproto.StackModeBelow)

		return nil
	}

	return fmt.Errorf("window not found: %s", title)
}

func (p *LinuxProvider) GetHandleByTitle(title string) string {
	if p.X == nil {
		return ""
	}
	clientIds, err := ewmh.ClientListGet(p.X)
	if err != nil {
		return ""
	}
	for _, id := range clientIds {
		name, _ := ewmh.WmNameGet(p.X, id)
		if name == "" {
			name, _ = icccm.WmNameGet(p.X, id)
		}
		if strings.Contains(name, title) {
			return fmt.Sprintf("%d", id)
		}
	}
	return ""
}

func (p *LinuxProvider) StackAbove(handleStr string, siblingStr string) error {
	if p.X == nil {
		return fmt.Errorf("X connection not available")
	}

	var id, sid uint32
	fmt.Sscanf(handleStr, "%d", &id)
	fmt.Sscanf(siblingStr, "%d", &sid)

	win := xwindow.New(p.X, xproto.Window(id))
	// Reconfigure window stacking relative to sibling
	win.StackSibling(xproto.Window(sid), xproto.StackModeAbove)
	return nil
}
