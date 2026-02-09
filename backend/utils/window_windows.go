//go:build windows

package utils

import (
	"fmt"
	"strings"
	"syscall"
	"unsafe"
)

var (
	user32                  = syscall.NewLazyDLL("user32.dll")
	procEnumWindows         = user32.NewProc("EnumWindows")
	procGetWindowTextW      = user32.NewProc("GetWindowTextW")
	procGetWindowRect       = user32.NewProc("GetWindowRect")
	procGetClientRect       = user32.NewProc("GetClientRect")
	procIsWindowVisible     = user32.NewProc("IsWindowVisible")
	procIsIconic            = user32.NewProc("IsIconic")
	procGetForegroundWindow = user32.NewProc("GetForegroundWindow")
	procSetWindowPos        = user32.NewProc("SetWindowPos")
	procGetWindow           = user32.NewProc("GetWindow")

	dwmapi                    = syscall.NewLazyDLL("dwmapi.dll")
	procDwmGetWindowAttribute = dwmapi.NewProc("DwmGetWindowAttribute")
)

const (
	DWMWA_EXTENDED_FRAME_BOUNDS = 9
	GW_HWNDPREV                 = 3
)

type winRect struct {
	Left   int32
	Top    int32
	Right  int32
	Bottom int32
}

type enumFindByTitleCtx struct {
	title          string
	requireVisible bool
	matched        syscall.Handle
	matchedTitle   string
}

type enumCollectWindowsCtx struct {
	windows *[]WindowInfo
}

func enumFindByTitleProc(hwnd syscall.Handle, lparam uintptr) uintptr {
	ctx := (*enumFindByTitleCtx)(unsafe.Pointer(lparam))
	if ctx == nil {
		return 0
	}

	if ctx.requireVisible {
		if isVisible, _, _ := procIsWindowVisible.Call(uintptr(hwnd)); isVisible == 0 {
			return 1
		}
	}

	var buf [256]uint16
	len, _, _ := procGetWindowTextW.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&buf[0])), 256)
	if len == 0 {
		return 1
	}

	titleStr := syscall.UTF16ToString(buf[:])
	if !strings.Contains(titleStr, ctx.title) {
		return 1
	}

	ctx.matched = hwnd
	ctx.matchedTitle = titleStr
	return 0
}

func enumCollectWindowsProc(hwnd syscall.Handle, lparam uintptr) uintptr {
	ctx := (*enumCollectWindowsCtx)(unsafe.Pointer(lparam))
	if ctx == nil || ctx.windows == nil {
		return 0
	}

	if isVisible, _, _ := procIsWindowVisible.Call(uintptr(hwnd)); isVisible == 0 {
		return 1
	}

	var title [256]uint16
	len, _, _ := procGetWindowTextW.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&title[0])), 256)
	if len == 0 {
		return 1
	}

	titleStr := syscall.UTF16ToString(title[:])

	var r winRect
	ret, _, _ := procDwmGetWindowAttribute.Call(
		uintptr(hwnd),
		uintptr(DWMWA_EXTENDED_FRAME_BOUNDS),
		uintptr(unsafe.Pointer(&r)),
		uintptr(unsafe.Sizeof(r)),
	)
	if ret != 0 {
		procGetWindowRect.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&r)))
	}

	isIconic, _, _ := procIsIconic.Call(uintptr(hwnd))

	*ctx.windows = append(*ctx.windows, WindowInfo{
		Handle:  fmt.Sprintf("%d", hwnd),
		Title:   titleStr,
		Process: "", // TODO: Get process name
		Rect: Rect{
			Left:   int(r.Left),
			Top:    int(r.Top),
			Right:  int(r.Right),
			Bottom: int(r.Bottom),
		},
		IsIconic: isIconic != 0,
	})

	return 1
}

var (
	enumFindByTitleCallback   = syscall.NewCallback(enumFindByTitleProc)
	enumCollectWindowsCallback = syscall.NewCallback(enumCollectWindowsProc)
)

// GetDWMFrameOffsetsByTitle finds a window by title and returns the pixel offsets
// between GetWindowRect (includes shadow) and DWM Extended Frame Bounds (visual bounds).
func GetDWMFrameOffsetsByTitle(title string) (FrameOffsets, error) {
	matched := findWindowByTitle(title)
	if matched == 0 {
		return FrameOffsets{}, fmt.Errorf("window not found: %s", title)
	}

	var wr winRect
	ret, _, _ := procGetWindowRect.Call(uintptr(matched), uintptr(unsafe.Pointer(&wr)))
	if ret == 0 {
		return FrameOffsets{}, fmt.Errorf("failed to get window rect")
	}

	var dwm winRect
	ret, _, _ = procDwmGetWindowAttribute.Call(
		uintptr(matched),
		uintptr(DWMWA_EXTENDED_FRAME_BOUNDS),
		uintptr(unsafe.Pointer(&dwm)),
		uintptr(unsafe.Sizeof(dwm)),
	)
	if ret != 0 {
		// DWM not available, offsets are zero (no shadow)
		return FrameOffsets{}, nil
	}

	return FrameOffsets{
		Top:    int(dwm.Top - wr.Top),
		Bottom: int(wr.Bottom - dwm.Bottom),
		Left:   int(dwm.Left - wr.Left),
		Right:  int(wr.Right - dwm.Right),
	}, nil
}

// GetWindowPhysicalWidthByTitle returns the physical pixel width of a window
// (from GetWindowRect, not DPI-scaled).
func GetWindowPhysicalWidthByTitle(title string) (int, error) {
	matched := findWindowByTitle(title)
	if matched == 0 {
		return 0, fmt.Errorf("window not found: %s", title)
	}

	var wr winRect
	ret, _, _ := procGetWindowRect.Call(uintptr(matched), uintptr(unsafe.Pointer(&wr)))
	if ret == 0 {
		return 0, fmt.Errorf("failed to get window rect")
	}

	return int(wr.Right - wr.Left), nil
}

// findWindowByTitle searches for a visible window containing the given title.
func findWindowByTitle(title string) syscall.Handle {
	ctx := enumFindByTitleCtx{title: title, requireVisible: true}
	procEnumWindows.Call(enumFindByTitleCallback, uintptr(unsafe.Pointer(&ctx)))
	return ctx.matched
}

func GetWindowDecorationHeightByTitle(title string) (int, error) {
	ctx := enumFindByTitleCtx{title: title, requireVisible: true}
	procEnumWindows.Call(enumFindByTitleCallback, uintptr(unsafe.Pointer(&ctx)))
	if ctx.matched == 0 {
		return 0, fmt.Errorf("window not found: %s", title)
	}

	var wr winRect
	ret, _, _ := procGetWindowRect.Call(uintptr(ctx.matched), uintptr(unsafe.Pointer(&wr)))
	if ret == 0 {
		return 0, fmt.Errorf("failed to get window rect: %s", ctx.matchedTitle)
	}

	var cr winRect
	ret, _, _ = procGetClientRect.Call(uintptr(ctx.matched), uintptr(unsafe.Pointer(&cr)))
	if ret == 0 {
		return 0, fmt.Errorf("failed to get client rect: %s", ctx.matchedTitle)
	}

	windowHeight := int(wr.Bottom - wr.Top)
	clientHeight := int(cr.Bottom - cr.Top)
	decorationHeight := windowHeight - clientHeight
	if decorationHeight < 0 {
		decorationHeight = 0
	}
	return decorationHeight, nil
}

type WindowsProvider struct{}

func NewWindowProvider() WindowProvider {
	return &WindowsProvider{}
}

func (p *WindowsProvider) Close() error {
	return nil
}

func (p *WindowsProvider) GetWindows() ([]WindowInfo, error) {
	var windows []WindowInfo
	ctx := enumCollectWindowsCtx{windows: &windows}
	procEnumWindows.Call(enumCollectWindowsCallback, uintptr(unsafe.Pointer(&ctx)))

	return windows, nil
}

func (p *WindowsProvider) GetWindowRect(handleStr string) (Rect, bool, error) {
	var handle uintptr
	fmt.Sscanf(handleStr, "%d", &handle)

	// Use DWM attribute for visual rect, fallback to GetWindowRect
	var r winRect
	ret, _, _ := procDwmGetWindowAttribute.Call(
		handle,
		uintptr(DWMWA_EXTENDED_FRAME_BOUNDS),
		uintptr(unsafe.Pointer(&r)),
		uintptr(unsafe.Sizeof(r)),
	)

	// DwmGetWindowAttribute returns 0 (S_OK) on success
	if ret != 0 {
		// Fallback to physical rect
		ret, _, _ = procGetWindowRect.Call(handle, uintptr(unsafe.Pointer(&r)))
		if ret == 0 {
			return Rect{}, false, fmt.Errorf("failed to get window rect")
		}
	}

	isIconic, _, _ := procIsIconic.Call(handle)

	return Rect{
		Left:   int(r.Left),
		Top:    int(r.Top),
		Right:  int(r.Right),
		Bottom: int(r.Bottom),
	}, isIconic != 0, nil
}

func (p *WindowsProvider) GetForegroundHandle() string {
	hwnd, _, _ := procGetForegroundWindow.Call()
	if hwnd == 0 {
		return ""
	}
	return fmt.Sprintf("%d", hwnd)
}

func ForceRaiseWindowByTitle(title string) error {
	ctx := enumFindByTitleCtx{title: title, requireVisible: true}
	procEnumWindows.Call(enumFindByTitleCallback, uintptr(unsafe.Pointer(&ctx)))
	if ctx.matched == 0 {
		return fmt.Errorf("window not found: %s", title)
	}

	const (
		SWP_NOSIZE     = 0x0001
		SWP_NOMOVE     = 0x0002
		SWP_SHOWWINDOW = 0x0040
	)

	// HWND_TOPMOST = -1
	procSetWindowPos.Call(uintptr(ctx.matched), ^uintptr(0), 0, 0, 0, 0, SWP_NOSIZE|SWP_NOMOVE|SWP_SHOWWINDOW)

	return nil
}

func ForceLowerWindowByTitle(title string) error {
	ctx := enumFindByTitleCtx{title: title, requireVisible: true}
	procEnumWindows.Call(enumFindByTitleCallback, uintptr(unsafe.Pointer(&ctx)))
	if ctx.matched == 0 {
		return fmt.Errorf("window not found: %s", title)
	}

	const (
		SWP_NOSIZE     = 0x0001
		SWP_NOMOVE     = 0x0002
		SWP_SHOWWINDOW = 0x0040
	)

	// HWND_NOTOPMOST = -2
	procSetWindowPos.Call(uintptr(ctx.matched), ^uintptr(1), 0, 0, 0, 0, SWP_NOSIZE|SWP_NOMOVE|SWP_SHOWWINDOW)

	return nil
}

func (p *WindowsProvider) GetHandleByTitle(title string) string {
	ctx := enumFindByTitleCtx{title: title, requireVisible: false}
	procEnumWindows.Call(enumFindByTitleCallback, uintptr(unsafe.Pointer(&ctx)))
	if ctx.matched != 0 {
		return fmt.Sprintf("%d", ctx.matched)
	}
	return ""
}

func (p *WindowsProvider) StackAbove(handleStr string, siblingStr string) error {
	var hwnd, hsibling uintptr
	fmt.Sscanf(handleStr, "%d", &hwnd)
	fmt.Sscanf(siblingStr, "%d", &hsibling)

	if hwnd == 0 || hsibling == 0 {
		return fmt.Errorf("invalid handles")
	}

	// HWND_TOP = 0
	// HWND_BOTTOM = 1
	// hWndInsertAfter: A handle to the window to precede the positioned window in the Z order.
	// To put A immediately above B, we need to put A after whatever is currently IN FRONT of B.
	prev, _, _ := procGetWindow.Call(hsibling, uintptr(GW_HWNDPREV))

	// SWP_NOSIZE | SWP_NOMOVE | SWP_NOACTIVATE = 0x01 | 0x02 | 0x10 = 0x13
	if prev == 0 {
		// Sibling is already at the top, so put handle at the top
		procSetWindowPos.Call(hwnd, 0, 0, 0, 0, 0, 0x0013)
	} else {
		procSetWindowPos.Call(hwnd, prev, 0, 0, 0, 0, 0x0013)
	}
	return nil
}

func ForceMoveResizeWindowByTitle(title string, x int, y int, width int, height int) error {
	ctx := enumFindByTitleCtx{title: title, requireVisible: true}
	procEnumWindows.Call(enumFindByTitleCallback, uintptr(unsafe.Pointer(&ctx)))
	if ctx.matched == 0 {
		return fmt.Errorf("window not found: %s", title)
	}

	return ForceMoveResizeWindow(fmt.Sprintf("%d", ctx.matched), x, y, width, height)
}

func parseHandleStr(handleStr string) (uintptr, error) {
	var handle uintptr
	if _, err := fmt.Sscanf(handleStr, "%d", &handle); err != nil {
		return 0, fmt.Errorf("invalid handle: %w", err)
	}
	if handle == 0 {
		return 0, fmt.Errorf("invalid handle")
	}
	return handle, nil
}

func GetDWMFrameOffsets(handleStr string) (FrameOffsets, error) {
	hwnd, err := parseHandleStr(handleStr)
	if err != nil {
		return FrameOffsets{}, err
	}

	var wr winRect
	ret, _, _ := procGetWindowRect.Call(hwnd, uintptr(unsafe.Pointer(&wr)))
	if ret == 0 {
		return FrameOffsets{}, fmt.Errorf("failed to get window rect")
	}

	var dwm winRect
	ret, _, _ = procDwmGetWindowAttribute.Call(
		hwnd,
		uintptr(DWMWA_EXTENDED_FRAME_BOUNDS),
		uintptr(unsafe.Pointer(&dwm)),
		uintptr(unsafe.Sizeof(dwm)),
	)
	if ret != 0 {
		return FrameOffsets{}, nil
	}

	return FrameOffsets{
		Top:    int(dwm.Top - wr.Top),
		Bottom: int(wr.Bottom - dwm.Bottom),
		Left:   int(dwm.Left - wr.Left),
		Right:  int(wr.Right - dwm.Right),
	}, nil
}

func GetWindowPhysicalWidth(handleStr string) (int, error) {
	hwnd, err := parseHandleStr(handleStr)
	if err != nil {
		return 0, err
	}

	var wr winRect
	ret, _, _ := procGetWindowRect.Call(hwnd, uintptr(unsafe.Pointer(&wr)))
	if ret == 0 {
		return 0, fmt.Errorf("failed to get window rect")
	}

	return int(wr.Right - wr.Left), nil
}

func ForceMoveResizeWindow(handleStr string, x int, y int, width int, height int) error {
	hwnd, err := parseHandleStr(handleStr)
	if err != nil {
		return err
	}

	const (
		SWP_NOZORDER   = 0x0004
		SWP_NOACTIVATE = 0x0010
	)

	procSetWindowPos.Call(hwnd, 0, uintptr(x), uintptr(y), uintptr(width), uintptr(height), SWP_NOZORDER|SWP_NOACTIVATE)
	return nil
}
