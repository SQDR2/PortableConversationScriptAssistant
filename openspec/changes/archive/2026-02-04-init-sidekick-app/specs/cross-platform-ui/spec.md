## ADDED Requirements

### Requirement: Mini Mode UI
The UI SHALL support a "Mini Mode" that consumes minimal screen real estate.

#### Scenario: App is idle
- **WHEN** the app is not in focus or explicitly set to Mini Mode
- **THEN** it displays as a thin vertical bar or small icon attached to the target window

### Requirement: Glassmorphism Effect
The application window background SHALL support a blur/transparency effect where supported by the OS.

#### Scenario: Expanded mode appearance
- **WHEN** the main window is visible
- **THEN** the background is semi-transparent with a blur effect, revealing contents behind it

### Requirement: Script List Virtualization
The script list UI SHALL support efficient rendering of at least 10,000 items.

#### Scenario: Scrolling large list
- **WHEN** user scrolls quickly through a list of 10,000 scripts
- **THEN** the UI remains responsive (60fps) without lag
