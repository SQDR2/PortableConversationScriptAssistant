## Context

The application currently serves as a conversation script assistant. Users need to translate text during interactions with foreign clients. We are integrating Tencent Cloud Text Translation (TMT) service through a new backend service and exposing it via a persistent UI component. The application uses Wails (Go backend + Vue frontend).

## Goals / Non-Goals

**Goals:**

- Implement `TranslationService` in Go to communicate with Tencent Cloud API.
- Create a configuration mechanism (`config.json`) for API credentials.
- Add a persistent "Translation Bar" in the frontend `MainLayout`.
- Support source/target language selection and real-time/manual translation.

**Non-Goals:**

- Offline translation.
- Text-to-Speech (TTS).
- OCR (Image translation).
- Complex UI for managing API keys (initially manual file editing).

## Decisions

### 1. Configuration Management

- **Decision**: Use a simple `config.json` file in the application execution directory.
- **Rationale**: Keeps implementation simple and allows users to easily provide keys without building a complex settings UI immediately. Secrets remain local.
- **Alternative**: Environment variables (harder for non-tech users) or Encrypted DB storage (overkill for v1).

### 2. UI Placement

- **Decision**: Fixed footer in `MainLayout`.
- **Rationale**: Translation is a tool used _alongside_ scripts. It should be always visible but not obstructing the main content.
- **Structure**:
  ```
  [ Language Selectors ]
  [ Input (adjustable height) ]
  [ Output (Read-only + Copy) ]
  ```

### 3. Backend Implementation

- **Decision**: Use `tencentcloud-sdk-go`.
- **Rationale**: Official SDK handles signing and requests robustly.
- **Service**: `TranslationService` bound to Wails context.

## Risks / Trade-offs

- **Risk**: API Key exposure.
  - _Mitigation_: `config.json` is local. User is responsible for their own key.
- **Risk**: Network Latency/Failure.
  - _Mitigation_: UI must show "Translating..." state and handle errors gracefully (e.g., "Network Error").
- **Trade-off**: `MainLayout` vertical space.
  - _Impact_: The footer will consume some screen real estate.
  - _Mitigation_: Keep it compact.
