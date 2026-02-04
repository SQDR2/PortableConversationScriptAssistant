## Context
This project aims to build a desktop "sidekick" application that enhances a user's workflow by providing context-aware script management. The application needs to run on Windows and Linux, tightly integrating with the window management system of the OS to "attach" itself to a target application (e.g., WeChat).

## Goals / Non-Goals

**Goals:**
- **Cross-Platform Architecture:** Use Wails (Go + Web) to share UI code while leveraging native OS APIs.
- **Window Following:** Implement a robust "magnetic" window behavior that follows the target app's position and visibility.
- **Performance:** Ensure minimal footprint (CPU/RAM) as it runs constantly in the background.
- **Search Latency:** Instant (<50ms) search results for scripts using FTS.

**Non-Goals:**
- **Mobile Support:** No Android/iOS capability planned.
- **Cloud Sync:** Local storage only for this version; no cloud backend.

## Decisions

### 1. Technology Stack: Wails
- **Decision:** Use Wails v2 over Electron or Flutter.
- **Rationale:** Wails produces significantly smaller binaries (native Go + WebView) compared to Electron. It offers easier access to low-level system APIs (via Go's  and ) which is critical for the "window following" feature compared to Flutter's FFI which can be more verbose for complex struct mapping.
- **Alternatives:**
  - *Flutter*: Good UI, but FFI for X11/Win32 is less ergonomic than Go.
  - *Electron*: Too heavy (Chromium bundle) for a "sidekick" tool that should be lightweight.

### 2. Database: SQLite with GORM
- **Decision:** Use SQLite embedded via CGO (or a pure Go driver if FTS isn't needed, but we need FTS).
- **Rationale:** We need "Fuzzy Keyword Search". SQLite's FTS5 extension provides robust full-text search capabilities out of the box. GORM provides a convenient ORM layer for the  model.

### 3. Window Tracking Strategy
- **Decision:** Polling mechanism in a Go goroutine.
- **Rationale:** Hooking global OS events (like  on Windows) is more efficient but complex and prone to flagging by AV software. A polling interval (e.g., 50ms) is simpler to implement, cross-platform friendly, and sufficient for visual sync without perceivable lag.
- **Linux Specific:** Use  (XLib) for window constraints. Wayland support is explicitly *out of scope* for V1 due to protocol restrictions.

## Risks / Trade-offs

- **Risk:** **Linux Wayland Support**
  - *Mitigation:* Explicitly state support for X11 only. Detect session type and warn user if on Wayland.
- **Risk:** **Windows Anti-Cheat/AV**
  - *Mitigation:* Ensure we only read window rects and do not inject code. Sign the binary if possible (future scope).
- **Risk:** **Polling Jitter**
  - *Mitigation:* Use a smoothing algorithm (Linear Interpolation - Lerp) if raw updates look jerky.

## Migration Plan
N/A - New Project.
