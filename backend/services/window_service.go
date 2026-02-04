package services

import (
	"context"
	"sidekick/backend/utils"
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
	return windows
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

	// Check if changed
	shouldForceUpdate := s.forceUpdate || s.initialAlignRetries > 0
	if s.hasLast && s.lastRect == rect && s.lastIconic == isIconic && !shouldForceUpdate {
		return
	}

	// Restore Skip Logic: When target is restored, WM needs time to update coordinates
	// If it was just restored, skip positioning for a few ticks to avoid jumping to (0,0)
	if s.hasLast && s.lastIconic && !isIconic {
		s.restoreSkipCount = 2 // Skip 100ms (2 * 50ms) for stability
		s.forceUpdate = true
		s.hasLast = false
		runtime.LogInfof(s.ctx, "Target restored, delaying positioning for stability")
	}

	if isIconic {
		if !s.isSidekickMinimized {
			runtime.WindowMinimise(s.ctx)
			s.isSidekickMinimized = true
		}
	} else {
		// Ensure window is not minimized or maximized before aligning
		runtime.WindowUnminimise(s.ctx)
		runtime.WindowUnmaximise(s.ctx)
		if s.isSidekickMinimized {
			s.isSidekickMinimized = false
		}
		if s.isSidekickHidden {
			runtime.WindowShow(s.ctx)
			s.isSidekickHidden = false
		}

		// Stability: Wait for WM to settle if recently restored
		if s.restoreSkipCount > 0 {
			s.restoreSkipCount--
			runtime.LogInfof(s.ctx, "Skipping position update (restore delay: %d ticks left)", s.restoreSkipCount)
			// Don't update lastRect to force re-check once settled
			s.lastIconic = isIconic // But update iconic state
			return
		}

		// Stick to the RIGHT side of the target (outside)
		if !s.decorationKnown {
			if dh, err := utils.GetWindowDecorationHeightByTitle("sidekick"); err == nil {
				s.decorationHeight = dh
				s.decorationKnown = true
				runtime.LogInfof(s.ctx, "Detected sidekick decoration height: %d", dh)
			} else {
				runtime.LogErrorf(s.ctx, "Failed to detect decoration height: %v", err)
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
		actualW, actualH := runtime.WindowGetSize(s.ctx)
		if actualX != newX || actualY != newY {
			err := utils.ForceMoveResizeWindowByTitle("sidekick", newX, newY, actualW, newHeight)
			if err != nil {
				runtime.LogErrorf(s.ctx, "Force move sidekick failed: %v", err)
			} else {
				actualX, actualY = runtime.WindowGetPosition(s.ctx)
				actualW, actualH = runtime.WindowGetSize(s.ctx)
			}
		}
		runtime.LogInfof(s.ctx, "Aligned sidekick to target: x=%d y=%d h=%d (actual: x=%d y=%d w=%d h=%d)", newX, newY, newHeight, actualX, actualY, actualW, actualH)
		if s.initialAlignRetries > 0 {
			s.initialAlignRetries--
		}
		if s.initialAlignRetries == 0 {
			s.forceUpdate = false
		}
	}

	s.lastRect = rect
	s.lastIconic = isIconic
	s.hasLast = true

	// Emit event to frontend for UI updates if needed
	runtime.EventsEmit(s.ctx, "window-position-update", map[string]interface{}{
		"rect":     rect,
		"isIconic": isIconic,
	})
}

// RuntimeWindowIsHidden removed as it is now tracked in the service instance
