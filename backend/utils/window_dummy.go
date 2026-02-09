//go:build !windows && !linux && !darwin

package utils

import "fmt"

type DummyProvider struct{}

func NewWindowProvider() WindowProvider {
	return &DummyProvider{}
}

func (p *DummyProvider) Close() error {
	return nil
}

func (p *DummyProvider) GetWindows() ([]WindowInfo, error) {
	return nil, fmt.Errorf("platform not supported")
}

func (p *DummyProvider) GetWindowRect(handle string) (Rect, bool, error) {
	return Rect{}, false, fmt.Errorf("platform not supported")
}

func (p *DummyProvider) GetForegroundHandle() string {
	return ""
}

func ForceRaiseWindowByTitle(title string) error {
	return fmt.Errorf("force raise not supported on this platform")
}

func ForceLowerWindowByTitle(title string) error {
	return fmt.Errorf("force lower not supported on this platform")
}

func (p *DummyProvider) GetHandleByTitle(title string) string {
	return ""
}

func (p *DummyProvider) StackAbove(handle string, sibling string) error {
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
