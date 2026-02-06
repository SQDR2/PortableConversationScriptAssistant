package services

import (
	"context"
	"sidekick/backend/utils"
	"strings"
	"sync"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type WindowService struct {
	ctx      context.Context
	provider utils.WindowProvider
	ticker   *time.Ticker
	stopChan chan struct{}

	mu         sync.Mutex
	targetHWID string

	// Optimization: Cache last state
	lastRect   utils.Rect
	lastIconic bool
	hasLast    bool

	// Stability: Track local visibility and restoration delay
	isSidekickHidden    bool
	isSidekickMinimized bool
	restoreSkipCount    int
	forceUpdate         bool
	initialAlignRetries int
	decorationHeight    int
	decorationKnown     bool
	sidekickHWID        string
	lastTargetFocused   bool
	isAlwaysOnTop       bool
}

func NewWindowService() *WindowService {
	return &WindowService{
		provider: utils.NewWindowProvider(),
		stopChan: make(chan struct{}),
	}
}

func (s *WindowService) Startup(ctx context.Context) {
	s.ctx = ctx
	s.StartPolling()
}

func (s *WindowService) Shutdown(ctx context.Context) {
	s.StopPolling()
	if s.provider != nil {
		s.provider.Close()
	}
}

func (s *WindowService) GetStartApps() []utils.WindowInfo {
	windows, err := s.provider.GetWindows()
	if err != nil {
		runtime.LogErrorf(s.ctx, "Failed to get windows: %v", err)
		return []utils.WindowInfo{}
	}

	// Initialize with empty slice to ensure it returns [] instead of null in JSON
	filtered := make([]utils.WindowInfo, 0)
	for _, w := range windows {
		// Filter out own window
		if strings.EqualFold(w.Title, "sidekick") {
			continue
		}
		filtered = append(filtered, w)
	}

	return filtered
}

func (s *WindowService) SetTarget(handle string) {
	s.mu.Lock()
	s.targetHWID = handle
	s.hasLast = false // Force immediate position update
	s.forceUpdate = true
	s.restoreSkipCount = 0
	s.initialAlignRetries = 3
	s.isSidekickMinimized = false
	s.decorationKnown = false
	s.mu.Unlock()

	// Cleanup any lingering sticky states
	runtime.WindowSetAlwaysOnTop(s.ctx, false)
	utils.ForceLowerWindowByTitle("sidekick")

	runtime.LogInfof(s.ctx, "Target set to: %s", handle)

	// Immediately align to new target
	s.checkTarget()
}

func (s *WindowService) StartPolling() {
	s.ticker = time.NewTicker(50 * time.Millisecond) // 50ms interval as per spec
	go func() {
		for {
			select {
			case <-s.ticker.C:
				s.checkTarget()
			case <-s.stopChan:
				return
			}
		}
	}()
}

func (s *WindowService) StopPolling() {
	if s.ticker != nil {
		s.ticker.Stop()
	}
	// Non-blocking send or close
	select {
	case s.stopChan <- struct{}{}:
	default:
	}
}

func (s *WindowService) checkTarget() {
	s.mu.Lock()
	target := s.targetHWID
	s.mu.Unlock()

	if target == "" {
		s.hasLast = false
		return
	}

	rect, isIconic, err := s.provider.GetWindowRect(target)
	if err != nil {
		s.hasLast = false
		runtime.LogErrorf(s.ctx, "Failed to get target rect: %v", err)
		// Window might be closed or lost
		return
	}

	foreground := s.provider.GetForegroundHandle()
	isTargetForeground := (foreground != "" && foreground == target)

	// 1. Visibility Synchronizer (Min/Max/Restore)
	if isIconic {
		if !s.isSidekickMinimized {
			runtime.WindowMinimise(s.ctx)
			s.isSidekickMinimized = true
			s.isSidekickHidden = true
			s.forceUpdate = false
			s.initialAlignRetries = 0
			runtime.LogInfof(s.ctx, "Target minimized, sidekick following")
		}
		// In iconic state, we update tracker and STOP early to avoid phantom popups
		goto update_tracker
	}

	if s.isSidekickMinimized {
		// Target RESTORED
		runtime.WindowUnminimise(s.ctx)
		runtime.WindowShow(s.ctx)
		s.isSidekickMinimized = false
		s.isSidekickHidden = false
		s.forceUpdate = true
		s.initialAlignRetries = 3
		runtime.LogInfof(s.ctx, "Target restored, resuming sidekick")
	}

	if s.sidekickHWID == "" {
		s.sidekickHWID = s.provider.GetHandleByTitle("sidekick")
	}

	// 2. Relative Stacking & Visibility (Active Guard)
	if s.isSidekickHidden {
		runtime.WindowShow(s.ctx)
		s.isSidekickHidden = false
	}

	if s.sidekickHWID != "" && target != "" {
		_ = s.provider.StackAbove(s.sidekickHWID, target)
	}

	// Stability: Skip specific alignment logic during restore transient period
	if s.restoreSkipCount > 0 {
		s.restoreSkipCount--
		goto update_tracker
	}

	// 3. Alignment Logic
	{
		// Stick to the RIGHT side of the target (outside)
		if !s.decorationKnown {
			if dh, err := utils.GetWindowDecorationHeightByTitle("sidekick"); err == nil {
				s.decorationHeight = dh
				s.decorationKnown = true
			}
		}
		sw, _ := runtime.WindowGetSize(s.ctx)
		newHeight := rect.Bottom - rect.Top
		if s.decorationKnown && newHeight > s.decorationHeight {
			newHeight = newHeight - s.decorationHeight
		}
		if newHeight > 0 {
			runtime.WindowSetSize(s.ctx, sw, newHeight)
		}
		newX := rect.Right
		newY := rect.Top

		runtime.WindowSetPosition(s.ctx, newX, newY)
		actualX, actualY := runtime.WindowGetPosition(s.ctx)
		if actualX != newX || actualY != newY {
			w, _ := runtime.WindowGetSize(s.ctx)
			_ = utils.ForceMoveResizeWindowByTitle("sidekick", newX, newY, w, newHeight)
		}

		if s.initialAlignRetries > 0 {
			s.initialAlignRetries--
		} else {
			s.forceUpdate = false
		}
	}

update_tracker:
	s.lastRect = rect
	s.lastIconic = isIconic
	s.lastTargetFocused = isTargetForeground
	s.hasLast = true

	// Emit event to frontend
	runtime.EventsEmit(s.ctx, "window-position-update", map[string]interface{}{
		"rect":     rect,
		"isIconic": isIconic,
	})
}

// RuntimeWindowIsHidden removed as it is now tracked in the service instance
