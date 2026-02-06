## ADDED Requirements

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
