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

### Requirement: Credential Prompt

The application SHALL check the translation configuration status on startup. If configuration is missing, it MUST prompt the user to input the required credentials before allowing access to translation features.

#### Scenario: Missing configuration on startup

- **WHEN** the main application layout is mounted
- **AND** the backend reports that translation configuration is missing
- **THEN** a modal configuration dialog is displayed
- **AND** the dialog prevents interaction with the rest of the application (or at least translation features) until resolved

#### Scenario: Successful configuration entry

- **WHEN** the user enters valid credentials and attempts to save
- **THEN** the application sends the credentials to the backend
- **AND** upon success, the dialog closes
- **AND** translation features become active

### Requirement: Configuration Dialog UI

The configuration dialog SHALL contain input fields for `SecretId` and `SecretKey`, and a mechanism to submit them.
All UI text MUST be in Simplified Chinese.

#### Scenario: Dialog content

- **WHEN** the dialog is displayed
- **THEN** it shows two text input fields (SecretId, SecretKey)
- **AND** a "Save" or "Confirm" button
- **AND** all labels and messages are in Chinese

### Requirement: Manual Configuration Access

The application SHALL provide an option in the main navigation drawer to manually open the configuration dialog, allowing users to update credentials at any time.

#### Scenario: Open from Drawer

- **WHEN** user opens the left side drawer
- **AND** clicks the "Translation Config" (translated) item
- **THEN** the configuration dialog opens
