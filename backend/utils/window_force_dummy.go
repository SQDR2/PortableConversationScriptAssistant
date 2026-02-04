//go:build !linux

package utils

import "fmt"

func ForceMoveResizeWindowByTitle(title string, x int, y int, width int, height int) error {
	return fmt.Errorf("force move not supported on this platform")
}
