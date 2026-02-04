//go:build !windows && !linux

package utils

import "fmt"

func GetWindowDecorationHeightByTitle(title string) (int, error) {
	return 0, fmt.Errorf("decoration height not supported on this platform")
}
