## Why

Users communicating across multiple platforms (WeChat, DingTalk, etc.) need a "sidekick" assistant that manages conversation scripts and follows their active window, appearing and disappearing contextually to provide instant access to pre-written text without switching focus.

## What Changes

- Initialize a new Wails (Go + Vue + Quasar) project structure.
- Implement a "sidekick" window mode that attaches to a specific user-selected target process/window.
- Create a local SQLite database engine for storing scripts.
- Build a polished "Glassmorphism" UI for script management (CRUD) and search.
- Support importing scripts from  files with custom delimiters.
- Support fuzzy keyword search for scripts.

## Capabilities

### New Capabilities

- `window-tracking`: Core logic for detecting, following, and mimicking the visibility state of a target window on Windows and Linux.
- `script-management`: Database schema and logic for creating, reading, updating, deleting, and importing conversation scripts.
- `cross-platform-ui`: The specific UI implementation using Quasar, including the "Mini Mode" and "Expanded Mode" with glassmorphism effects.

### Modified Capabilities

- None.

## Impact

- **New Project**: This change establishes the entire codebase for the new application.
- **Dependencies**: Introduces Wails, Go, Vue 3, Quasar, and GORM/SQLite.
- **System**: Requires system-level access (User32.dll on Windows, X11/EWMH on Linux) for window management.
