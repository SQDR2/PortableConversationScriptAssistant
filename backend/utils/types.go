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

type WindowProvider interface {
	GetWindows() ([]WindowInfo, error)
	GetWindowRect(handle string) (Rect, bool, error)
	GetForegroundHandle() string
	GetHandleByTitle(title string) string
	StackAbove(handle string, siblingHandle string) error
	Close() error
}
