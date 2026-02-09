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
	frameOffsets        utils.FrameOffsets
	frameOffsetsKnown   bool
	sidekickHWID        string
	lastSidekickLookup  time.Time
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
	s.frameOffsetsKnown = false
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
	defer func() {
		if r := recover(); r != nil {
			utils.LogError("check_target_panic", r)
		}
	}()

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
		// Avoid hot-looping window enumeration if the handle isn't available yet.
		if s.lastSidekickLookup.IsZero() || time.Since(s.lastSidekickLookup) > time.Second {
			s.lastSidekickLookup = time.Now()
			s.sidekickHWID = s.provider.GetHandleByTitle("sidekick")
		}
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
		// Get sidekick's DWM frame offsets (shadow compensation) for precise alignment
		if !s.frameOffsetsKnown && s.sidekickHWID != "" {
			if offsets, err := utils.GetDWMFrameOffsets(s.sidekickHWID); err == nil {
				s.frameOffsets = offsets
				s.frameOffsetsKnown = true
			}
		}

		// rect is the target's DWM visual bounds (no shadow)
		targetVisualHeight := rect.Bottom - rect.Top

		// Calculate SetWindowPos parameters that align sidekick's visual bounds
		// with the target's visual bounds.
		// shadow offsets: the gap between GetWindowRect and DWM visual bounds.
		newHeight := targetVisualHeight + s.frameOffsets.Top + s.frameOffsets.Bottom
		newX := rect.Right - s.frameOffsets.Left // sidekick visual left edge = target visual right edge
		newY := rect.Top - s.frameOffsets.Top    // sidekick visual top = target visual top

		// Get current sidekick physical pixel width (bypass Wails DPI scaling)
		sw := 350
		if s.sidekickHWID != "" {
			w, err := utils.GetWindowPhysicalWidth(s.sidekickHWID)
			if err != nil {
				// Handle may have become invalid (e.g., window recreated); clear and re-resolve later.
				s.sidekickHWID = ""
				s.frameOffsetsKnown = false
			} else if w > 0 {
				sw = w
			}
		}

		// Directly use native APIs to bypass Wails DPI scaling and monitor offset
		if newHeight > 0 {
			if s.sidekickHWID != "" {
				_ = utils.ForceMoveResizeWindow(s.sidekickHWID, newX, newY, sw, newHeight)
			} else {
				_ = utils.ForceMoveResizeWindowByTitle("sidekick", newX, newY, sw, newHeight)
			}
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
