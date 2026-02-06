## ADDED Requirements

### Requirement: Query Configuration Status

The `config-management` capability SHALL expose a mechanism to query whether the required configuration (specifically, Tencent Cloud credentials) is valid and present.

#### Scenario: Config is invalid or missing

- **WHEN** the system configuration is loaded with empty SecretId or SecretKey
- **THEN** the status query returns `false` (not configured)

#### Scenario: Config is valid

- **WHEN** the system configuration contains non-empty SecretId and SecretKey
- **THEN** the status query returns `true` (configured)

### Requirement: Update and Persist Configuration

The capability SHALL provide a method to update the configuration credentials at runtime. This method MUST persist the new credentials to the underlying storage (e.g., `config.json`) and trigger a reload of dependent services without requiring an application restart.

#### Scenario: Update credentials

- **WHEN** `UpdateCredentials` is called with a new SecretId and SecretKey
- **THEN** the `config.json` file on disk is updated with these values
- **AND** the in-memory configuration is updated
- **AND** dependent services are notified or re-initialized
