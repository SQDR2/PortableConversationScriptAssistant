## ADDED Requirements

### Requirement: Display Translation Interface

The application SHALL display a translation bar fixed at the bottom of the main layout. It MUST contain:

- Language selection dropdowns for Source and Target languages.
- An input text area (defaulting to Source Language).
- A read-only output text area (defaulting to Target Language).

#### Scenario: Layout rendering

- **WHEN** user opens any page in the application
- **THEN** the translation bar is visible at the bottom

### Requirement: Execute Translation

The interface SHALL automatically trigger translation when the user stops typing in the input area for a specified duration (e.g., 500ms debounce), OR when the user explicitly actions it (e.g., Enter).

#### Scenario: Debounced Translation

- **WHEN** user types "Hello" and pauses
- **THEN** application calls backend translation service
- **AND** updates the output area with the result

### Requirement: Copy Result

The interface SHALL allow the user to easily copy the translated text to the clipboard.

#### Scenario: Click to copy

- **WHEN** user clicks on the output text area or a "Copy" button
- **THEN** the translated text is copied to the system clipboard
- **AND** a notification confirms the action
