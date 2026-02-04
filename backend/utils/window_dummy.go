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
