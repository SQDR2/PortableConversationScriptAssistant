//go:build windows

package utils

import (
	"fmt"
	"strings"
	"syscall"
	"unsafe"
)

var (
	user32              = syscall.NewLazyDLL("user32.dll")
	procEnumWindows     = user32.NewProc("EnumWindows")
	procGetWindowTextW  = user32.NewProc("GetWindowTextW")
	procGetWindowRect   = user32.NewProc("GetWindowRect")
	procGetClientRect   = user32.NewProc("GetClientRect")
	procIsWindowVisible = user32.NewProc("IsWindowVisible")
	procIsIconic        = user32.NewProc("IsIconic")
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

		var r winRect
		procGetWindowRect.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&r)))

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

	var r winRect
	ret, _, _ := procGetWindowRect.Call(handle, uintptr(unsafe.Pointer(&r)))
	if ret == 0 {
		return Rect{}, false, fmt.Errorf("failed to get window rect")
	}

	isIconic, _, _ := procIsIconic.Call(handle)

	return Rect{
		Left:   int(r.Left),
		Top:    int(r.Top),
		Right:  int(r.Right),
		Bottom: int(r.Bottom),
	}, isIconic != 0, nil
}
