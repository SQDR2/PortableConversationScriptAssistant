## 1. Backend Implementation

- [x] 1.1 Update `backend/utils/config.go` with `SaveConfig` function.
- [x] 1.2 Update `backend/services/translation_service.go` to add `GetConfigStatus` method (check for empty strings).
- [x] 1.3 Update `backend/services/translation_service.go` to add `UpdateCredentials` method (update struct, save to file, re-run init logic).
- [x] 1.4 Expose new methods to Wails runtime (ensure they are public and bound).

## 2. Frontend Implementation

- [x] 2.1 Update `frontend/wailsjs` bindings (running `wails dev` or manual generation).
- [x] 2.2 Create `frontend/src/components/ConfigDialog.vue` (or implement directly in MainLayout) using `q-dialog`, input fields, and save button.
- [x] 2.3 Modify `frontend/src/layouts/MainLayout.vue` to check `GetConfigStatus` on mount.
- [x] 2.4 Implement logic to show dialog if unconfigured, and handle save action.

## 3. Verification

- [x] 3.1 Verify startup with empty `config.json` triggers the prompt.
- [x] 3.2 Verify entering credentials saves to `config.json` correctly.
- [x] 3.3 Verify translation works immediately after saving without restart.
- [x] 3.4 Verify restarting app with valid config does NOT trigger prompt.

## 4. Refinement

- [x] 4.1 Localize `frontend/src/components/ConfigDialog.vue` to Simplified Chinese.
- [x] 4.2 Add "Translation Config" item to Drawer in `frontend/src/layouts/MainLayout.vue`.
- [x] 4.3 Link Drawer item to open the ConfigDialog.
