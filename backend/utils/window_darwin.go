//go:build darwin

package utils

import "fmt"

type DarwinProvider struct{}

func NewWindowProvider() WindowProvider {
	return &DarwinProvider{}
}

func (p *DarwinProvider) Close() error {
	return nil
}

func (p *DarwinProvider) GetWindows() ([]WindowInfo, error) {
	return nil, fmt.Errorf("platform not supported")
}

func (p *DarwinProvider) GetWindowRect(handle string) (Rect, bool, error) {
	return Rect{}, false, fmt.Errorf("platform not supported")
}

func (p *DarwinProvider) GetForegroundHandle() string {
	return ""
}

func ForceRaiseWindowByTitle(title string) error {
	return fmt.Errorf("force raise not supported on this platform")
}

func ForceLowerWindowByTitle(title string) error {
	return fmt.Errorf("force lower not supported on this platform")
}

func (p *DarwinProvider) GetHandleByTitle(title string) string {
	return ""
}

func (p *DarwinProvider) StackAbove(handle string, sibling string) error {
	return nil
}

func GetWindowDecorationHeightByTitle(title string) (int, error) {
	return 0, nil
}

func GetDWMFrameOffsetsByTitle(title string) (FrameOffsets, error) {
	return FrameOffsets{}, nil
}

func GetWindowPhysicalWidthByTitle(title string) (int, error) {
	return 0, nil
}

func ForceMoveResizeWindowByTitle(title string, x int, y int, width int, height int) error {
	return nil
}
