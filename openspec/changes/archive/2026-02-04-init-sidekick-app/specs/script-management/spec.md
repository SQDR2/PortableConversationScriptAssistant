## ADDED Requirements

### Requirement: Create and Edit Scripts
The system SHALL allow users to create, read, update, and delete conversation scripts.

#### Scenario: Create new script
- **WHEN** user clicks "New Script" and enters "Hello World"
- **THEN** a new script record is created in the database with the content "Hello World"

#### Scenario: Update script
- **WHEN** user edits an existing script content
- **THEN** the changes are saved to the database

### Requirement: Import Scripts from Text File
The system SHALL parse imported text files into individual script entries based on a delimiter.

#### Scenario: Import with line delimiter
- **WHEN** user imports a file containing "Line 1\n\nLine 2" and selects "Double Newline" as delimiter
- **THEN** the system creates 2 distinct script entries

### Requirement: Search Scripts
The system SHALL provide full-text search capabilities for scripts.

#### Scenario: Fuzzy search
- **WHEN** user types "refund" in the search bar
- **THEN** the list filters to show scripts containing "refund", "funds", etc. (depending on fuzzy logic configuration)
