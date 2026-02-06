## 1. Setup & Configuration

- [x] 1.1 Add Tencent Cloud SDK dependency to `go.mod`: `go get github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common` and `tencentcloud/tmt/v20180321`.
- [x] 1.2 Create `config.json` in root directory with template structure.
- [x] 1.3 Create `backend/utils/config.go` to define struct and load config from file.

## 2. Backend Implementation

- [x] 2.1 Create `backend/services/translation_service.go`.
- [x] 2.2 Implement `Startup` method to load config and initialize TMT client.
- [x] 2.3 Implement `Translate(text, source, target)` method wrapping TMT API.
- [x] 2.4 Register `TranslationService` in `main.go`.

## 3. Frontend Implementation

- [x] 3.1 Update `MainLayout.vue` to add a footer container.
- [x] 3.2 Implement UI components: Language selectors (Quasar `q-select`), Input/Output (`q-input`).
- [x] 3.3 Connect Frontend to Backend: Import `Translate` from wailsjs, implement call logic.
- [x] 3.4 Add "Copy" functionality and error handling (notifications).

## 4. Verification

- [x] 4.1 Verify connection to Tencent API with real/test key.
- [x] 4.2 Verify UI responsiveness and layout.

## 5. UI Refinement

- [x] 5.1 Set translation bar to hidden by default in `MainLayout.vue`.
- [x] 5.2 Verify `q-layout` view configuration (`lHh Lpr lFf`) ensures `lFf` (Footer) correctly compresses page content (Fixed via `style-fn` update).
- [ ] 5.3 Add transition animation for smoother toggling (Optional).

## 6. Feedback Fixes

- [x] 6.1 Fix Language Dropdown font color (currently incorrect/too light).
- [x] 6.2 Disable auto-debounce translation; trigger ONLY on Enter key.
- [x] 6.3 Verify and ensure "Click to Copy" functionality is obvious and working.

## 7. Round 2 Feedback Fixes

- [x] 7.1 Fix Language Dropdown **Popup Menu** font color (options are white-on-white).
- [x] 7.2 Prevent Output Box from shrinking/resizing when clicked/copied (Fixed height).

## 8. Round 3 Feedback Fixes

- [x] 8.1 Fix Output Box Width shrinking on click (force full width).

## 9. Round 4 Feedback Fixes

- [x] 9.1 (Retry) Align Target Language text to the right using CSS flex alignment.

## 10. Round 5 Feedback Fixes

- [x] 10.1 Set translation bar to **visible** by default.
