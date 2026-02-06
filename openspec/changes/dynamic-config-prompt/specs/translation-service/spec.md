## MODIFIED Requirements

### Requirement: Load Configuration

The service SHALL load translation API credentials (SecretId, SecretKey, Region, ProjectId) from a `config.json` file located in the application root directory. It MUST support dynamic reloading of these credentials when configuration is updated at runtime.

#### Scenario: Valid config file exists

- **WHEN** application starts
- **THEN** service successfully reads credentials from `config.json`

#### Scenario: Config file missing or invalid

- **WHEN** application starts and `config.json` is missing or credentials are empty
- **THEN** service enters an unconfigured state but does NOT crash
- **AND** `TranslationService` marks itself as unconfigured

## ADDED Requirements

### Requirement: Dynamic Re-initialization

The service SHALL provide a mechanism to re-initialize its internal TMT client using new credentials provided at runtime.

#### Scenario: Credential update

- **WHEN** new credentials are provided via `UpdateCredentials`
- **THEN** the service tears down the old client (if any)
- **AND** initializes a new TMT client with the new credentials
- **AND** updates its internal configuration state to "Configured"
