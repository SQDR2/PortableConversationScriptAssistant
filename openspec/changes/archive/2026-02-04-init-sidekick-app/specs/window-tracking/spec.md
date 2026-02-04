## ADDED Requirements

### Requirement: Select Target Process
The system SHALL provide a mechanism for the user to select a running application process to "attach" to.

#### Scenario: User selects a process from a list
- **WHEN** user opens the "Select Target" dialog
- **THEN** the system displays a list of visible application windows with their titles and icons
- **WHEN** user selects "WeChat.exe"
- **THEN** the system stores this process as the active target

### Requirement: Follow Target Position
The application window SHALL automatically update its position to remain adjacent to the target window when the target window is moved.

#### Scenario: Target window moves
- **WHEN** the target window (e.g., WeChat) is dragged to coordinate (X, Y)
- **THEN** the Sidekick app window moves to (X + TargetWidth, Y) (or user-configured relative offset) within 50ms

### Requirement: Sync Visibility State
The application window SHALL match the visibility state of the target window (Minimize, Restore, Hide).

#### Scenario: Target minimizes
- **WHEN** the user minimizes the target window to the taskbar
- **THEN** the Sidekick app window also hides/minimizes

#### Scenario: Target restores
- **WHEN** the user restores the target window
- **THEN** the Sidekick app window becomes visible again at the correct relative position
