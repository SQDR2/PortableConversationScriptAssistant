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
