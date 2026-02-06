## ADDED Requirements

### Requirement: Display Translation Interface

The application SHALL display a translation bar fixed at the bottom of the main layout. It MUST contain:

- Language selection dropdowns for Source and Target languages.
- An input text area (defaulting to Source Language).
- A read-only output text area (defaulting to Target Language).

#### Scenario: Layout rendering

- **WHEN** user opens application
- **THEN** the translation bar is visible by default

### Requirement: Toggle Visibility

The translation bar SHALL appear at the bottom when the user clicks the translation toggle button in the header. When visible, it MUST push the main content area up (compress scroll area) rather than overlaying it.

#### Scenario: Toggle On

- **WHEN** user clicks translation button
- **THEN** translation bar appears
- **AND** main script list height decreases to fit

#### Scenario: Toggle Off

- **WHEN** user clicks translation button again
- **THEN** translation bar disappears
- **AND** main script list height expands

### Requirement: Execute Translation

The interface SHALL trigger translation ONLY when the user explicitly actions it (e.g., pressing Enter). Automatic translation (debounce) SHALL be disabled to conserve API usage.

#### Scenario: Manual Translation

- **WHEN** user types "Hello" and presses Enter
- **THEN** application calls backend translation service
- **AND** updates the output area with the result

### Requirement: Copy Result

The interface SHALL allow the user to easily copy the translated text to the clipboard.

#### Scenario: Click to copy

- **WHEN** user clicks on the output text area or a "Copy" button
- **THEN** the translated text is copied to the system clipboard
- **AND** a notification confirms the action
