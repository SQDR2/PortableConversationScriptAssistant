## Why

The current application facilitates conversation scripts but lacks a tool for immediate cross-language communication. When interacting with foreign clients (e.g., on DingTalk, WeChat, WeChat Work), users often need to translate input text to a target language or verify received text. Adding an integrated translation tool will streamline this workflow, keeping the user context within the "Sidekick" assistant.

## What Changes

- **Backend**:
  - Add a new `TranslationService` in Go.
  - Integrate Tencent Cloud TMT (Text Machine Translation) SDK.
  - Implement a `config.json` loader to securely manage sensitive API keys (SecretId, SecretKey).
- **Frontend**:
  - Add a persistent "Translation Bar" at the bottom of the main layout.
  - Support real-time (debounced) or manual text translation.
  - Support language selection (Source/Target).
  - Provide "Copy" functionality for the translated result.
- **Config**:
  - Introduce `config.json` for backend configuration.

## Capabilities

### New Capabilities

- `translation-service`: Core translation logic interacting with Tencent Cloud API.
- `ui-translation-bar`: The visual component for user interaction (Input/Output/Language Select).

### Modified Capabilities

- None.

## Impact

- **New Dependencies**: `tencentcloud-sdk-go` for Go.
- **Files**:
  - `backend/services/translation_service.go`
  - `backend/config/config.go` (or similar)
  - `frontend/src/components/TranslationBar.vue` (or integration into MainLayout)
  - `config.json` (root directory)
