package utils

type Rect struct {
	Left   int `json:"left"`
	Top    int `json:"top"`
	Right  int `json:"right"`
	Bottom int `json:"bottom"`
}

type WindowInfo struct {
	Handle   string `json:"handle"`
	Title    string `json:"title"`
	Process  string `json:"process"`
	Rect     Rect   `json:"rect"`
	IsIconic bool   `json:"isIconic"`
}

// FrameOffsets represents the pixel difference between a window's
// GetWindowRect bounds and its DWM Extended Frame Bounds (visual bounds).
// On Windows 10/11, this captures the invisible shadow/border pixels.
type FrameOffsets struct {
	Top    int // dwm.Top - windowRect.Top (usually ≈ 0)
	Bottom int // windowRect.Bottom - dwm.Bottom (usually ≈ 7-8px)
	Left   int // dwm.Left - windowRect.Left
	Right  int // windowRect.Right - dwm.Right
}

type WindowProvider interface {
	GetWindows() ([]WindowInfo, error)
	GetWindowRect(handle string) (Rect, bool, error)
	GetForegroundHandle() string
	GetHandleByTitle(title string) string
	StackAbove(handle string, siblingHandle string) error
	Close() error
}
