## ADDED Requirements

### Requirement: Load Configuration

The service SHALL load translation API credentials (SecretId, SecretKey, Region, ProjectId) from a `config.json` file located in the application root directory.

#### Scenario: Valid config file exists

- **WHEN** application starts
- **THEN** service successfully reads credentials from `config.json`

#### Scenario: Config file missing

- **WHEN** application starts and `config.json` is missing
- **THEN** service handles the error gracefully (e.g., translation requests fail with specific error "Config missing")

### Requirement: Translate Text

The service SHALL provide a function to translate text using the Tencent Cloud TMT API. It MUST accept source text, source language, and target language as inputs.

#### Scenario: Successful translation

- **WHEN** client calls valid translate function with "Hello", source "auto", target "zh"
- **THEN** service returns "你好" (or equivalent) and no error

#### Scenario: API Error handling

- **WHEN** Tencent API returns an error (e.g., auth failure)
- **THEN** service returns the error message to the caller
