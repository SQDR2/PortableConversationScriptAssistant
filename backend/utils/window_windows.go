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

func GetWindowDecorationHeightByTitle(title string) (int, error) {
	var matched syscall.Handle
	var matchedTitle string

	cb := syscall.NewCallback(func(hwnd syscall.Handle, lparam uintptr) uintptr {
		if isVisible, _, _ := procIsWindowVisible.Call(uintptr(hwnd)); isVisible == 0 {
			return 1 // Continue
		}

		var buf [256]uint16
		len, _, _ := procGetWindowTextW.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&buf[0])), 256)
		if len == 0 {
			return 1
		}

		titleStr := syscall.UTF16ToString(buf[:])
		if !strings.Contains(titleStr, title) {
			return 1
		}

		matched = hwnd
		matchedTitle = titleStr
		return 0 // Stop enumeration
	})

	procEnumWindows.Call(cb, 0)

	if matched == 0 {
		return 0, fmt.Errorf("window not found: %s", title)
	}

	var wr winRect
	ret, _, _ := procGetWindowRect.Call(uintptr(matched), uintptr(unsafe.Pointer(&wr)))
	if ret == 0 {
		return 0, fmt.Errorf("failed to get window rect: %s", matchedTitle)
	}

	var cr winRect
	ret, _, _ = procGetClientRect.Call(uintptr(matched), uintptr(unsafe.Pointer(&cr)))
	if ret == 0 {
		return 0, fmt.Errorf("failed to get client rect: %s", matchedTitle)
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

	cb := syscall.NewCallback(func(hwnd syscall.Handle, lparam uintptr) uintptr {
		if isVisible, _, _ := procIsWindowVisible.Call(uintptr(hwnd)); isVisible == 0 {
			return 1 // Continue
		}

		var title [256]uint16
		len, _, _ := procGetWindowTextW.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&title[0])), 256)
		if len == 0 {
			return 1
		}

		titleStr := syscall.UTF16ToString(title[:])

		// Use DWM attribute for visual rect, fallback to GetWindowRect
		var r winRect
		ret, _, _ := procDwmGetWindowAttribute.Call(
			uintptr(hwnd),
			uintptr(DWMWA_EXTENDED_FRAME_BOUNDS),
			uintptr(unsafe.Pointer(&r)),
			uintptr(unsafe.Sizeof(r)),
		)
		if ret != 0 {
			// Fallback to physical rect
			procGetWindowRect.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&r)))
		}

		isIconic, _, _ := procIsIconic.Call(uintptr(hwnd))

		windows = append(windows, WindowInfo{
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
	})

	procEnumWindows.Call(cb, 0)

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
	var matched syscall.Handle

	cb := syscall.NewCallback(func(hwnd syscall.Handle, lparam uintptr) uintptr {
		if isVisible, _, _ := procIsWindowVisible.Call(uintptr(hwnd)); isVisible == 0 {
			return 1
		}

		var buf [256]uint16
		len, _, _ := procGetWindowTextW.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&buf[0])), 256)
		if len == 0 {
			return 1
		}

		titleStr := syscall.UTF16ToString(buf[:])
		if !strings.Contains(titleStr, title) {
			return 1
		}

		matched = hwnd
		return 0
	})

	procEnumWindows.Call(cb, 0)

	if matched == 0 {
		return fmt.Errorf("window not found: %s", title)
	}

	const (
		SWP_NOSIZE     = 0x0001
		SWP_NOMOVE     = 0x0002
		SWP_SHOWWINDOW = 0x0040
	)

	// HWND_TOPMOST = -1
	procSetWindowPos.Call(uintptr(matched), ^uintptr(0), 0, 0, 0, 0, SWP_NOSIZE|SWP_NOMOVE|SWP_SHOWWINDOW)

	return nil
}

func ForceLowerWindowByTitle(title string) error {
	var matched syscall.Handle

	cb := syscall.NewCallback(func(hwnd syscall.Handle, lparam uintptr) uintptr {
		if isVisible, _, _ := procIsWindowVisible.Call(uintptr(hwnd)); isVisible == 0 {
			return 1
		}

		var buf [256]uint16
		len, _, _ := procGetWindowTextW.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&buf[0])), 256)
		if len == 0 {
			return 1
		}

		titleStr := syscall.UTF16ToString(buf[:])
		if !strings.Contains(titleStr, title) {
			return 1
		}

		matched = hwnd
		return 0
	})

	procEnumWindows.Call(cb, 0)

	if matched == 0 {
		return fmt.Errorf("window not found: %s", title)
	}

	const (
		SWP_NOSIZE     = 0x0001
		SWP_NOMOVE     = 0x0002
		SWP_SHOWWINDOW = 0x0040
	)

	// HWND_NOTOPMOST = -2
	procSetWindowPos.Call(uintptr(matched), ^uintptr(1), 0, 0, 0, 0, SWP_NOSIZE|SWP_NOMOVE|SWP_SHOWWINDOW)

	return nil
}

func (p *WindowsProvider) GetHandleByTitle(title string) string {
	var matched syscall.Handle
	cb := syscall.NewCallback(func(hwnd syscall.Handle, lparam uintptr) uintptr {
		var buf [256]uint16
		len, _, _ := procGetWindowTextW.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&buf[0])), 256)
		if len != 0 {
			titleStr := syscall.UTF16ToString(buf[:])
			if strings.Contains(titleStr, title) {
				matched = hwnd
				return 0
			}
		}
		return 1
	})
	procEnumWindows.Call(cb, 0)
	if matched != 0 {
		return fmt.Sprintf("%d", matched)
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
	var matched syscall.Handle

	cb := syscall.NewCallback(func(hwnd syscall.Handle, lparam uintptr) uintptr {
		if isVisible, _, _ := procIsWindowVisible.Call(uintptr(hwnd)); isVisible == 0 {
			return 1
		}

		var buf [256]uint16
		len, _, _ := procGetWindowTextW.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&buf[0])), 256)
		if len == 0 {
			return 1
		}

		titleStr := syscall.UTF16ToString(buf[:])
		if !strings.Contains(titleStr, title) {
			return 1
		}

		matched = hwnd
		return 0
	})

	procEnumWindows.Call(cb, 0)

	if matched == 0 {
		return fmt.Errorf("window not found: %s", title)
	}

	const (
		SWP_NOZORDER   = 0x0004
		SWP_NOACTIVATE = 0x0010
	)

	procSetWindowPos.Call(uintptr(matched), 0, uintptr(x), uintptr(y), uintptr(width), uintptr(height), SWP_NOZORDER|SWP_NOACTIVATE)
	return nil
}
