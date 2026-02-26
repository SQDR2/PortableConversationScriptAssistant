# AGENTS.md — Sidekick (PortableConversationScriptAssistant)

## Project Overview
Sidekick is a Wails v2 desktop application (Go backend + Vue 3/Quasar frontend) designed for conversation script management.
Key features include a SQLite (GORM) database, window pinning (always-on-top), and Tencent Cloud TMT translation.

## Build, Lint, & Test Commands

### Development & Build
- **Dev Mode (Hot Reload):** `wails dev` (recompiles Go and runs Vite dev server)
- **Build App:** `wails build`
- **Cross-compile for Linux:** `wails build -platform linux/amd64`
- **Cross-compile for Windows:** `wails build -platform windows/amd64`
- **Frontend Only Dev:** `cd frontend && pnpm dev`
- **Frontend Build:** `cd frontend && pnpm build` (runs `vue-tsc --noEmit && vite build`)
- **Go Compile Check:** `go build ./...`
- **Go Static Analysis (Lint):** `go vet ./...`

### Testing (Go Backend)
- **Run all tests:** `go test ./...`
- **Run service tests:** `go test ./backend/services/...`
- **Run a single test by name:** `go test -run TestDbNilGuard ./backend/services/`
- **Run a single test verbosely:** `go test -v -run TestName ./path/to/package/`
> Note: Test coverage is currently minimal. No frontend tests exist. Always add Go tests when adding new backend features or fixing bugs.

## Architecture & Code Structure
- `main.go` / `app.go`: Application entry point, embeds frontend build, wires Wails services.
- `backend/db/`: Global `*gorm.DB` singleton (SQLite in WAL mode).
- `backend/models/`: GORM structs mapping to database tables (`Script`, `Category`).
- `backend/services/`: Core application logic (e.g., `ScriptService`, `WindowService`).
- `backend/utils/`: Configuration, logging, types, and platform-specific window implementation files.
- `frontend/src/`: Vue 3 + Quasar setup.
  - `pages/`, `components/`, `layouts/`, `router/`
  - `wailsjs/`: Auto-generated Wails Go bindings (**DO NOT EDIT**).

## Code Style — Go Backend

### Naming Conventions
- **Exported Identifiers:** `PascalCase` (e.g., `ScriptService`, `GetAllScripts`)
- **Unexported Identifiers:** `camelCase` (e.g., `logInfof`, `checkTarget`)
- **Constructors:** Use the `NewXxxService()` pattern.
- **Receivers:** Use short lowercase letters (e.g., `func (s *ScriptService) Method()`).

### Imports & Formatting
- Group imports into three separated blocks: 
  1. Standard library (`fmt`, `log`)
  2. Internal packages (`sidekick/backend/db`)
  3. Third-party packages (`gorm.io/gorm`)
- Use standard `gofmt` for all formatting.

### Types & Services
- Every Wails service must hold context for runtime operations:
  ```go
  type XxxService struct { ctx context.Context }
  func NewXxxService() *XxxService { return &XxxService{} }
  func (s *XxxService) Startup(ctx context.Context) { s.ctx = ctx }
  ```
- Protect shared mutable state with `sync.Mutex`. Always use `defer s.mu.Unlock()` immediately after `s.mu.Lock()`.

### Database Access & Error Handling
- **Crucial:** Always nil-check the global `db.DB` before use to prevent panics.
- Use `gorm.DeletedAt` for soft deletes.
- Wrap multi-step or dependent operations in `db.DB.Transaction(func(tx *gorm.DB) error { ... })`.
- Return `(value, error)` tuples from all service methods. Avoid deep nesting by returning early on error.
- Wrap errors with descriptive context: `fmt.Errorf("failed to fetch script: %w", err)`.
- Log errors using Wails runtime: `runtime.LogErrorf(s.ctx, "message: %v", err)`.

## Code Style — Vue 3 / TypeScript Frontend

### SFC Structure
Follow the standard `<script setup>` -> `<template>` -> `<style>` order:
```vue
<script setup lang="ts">
// Vue logic, imports, and composables
</script>

<template>
  <!-- Quasar UI components -->
</template>

<style lang="scss" scoped>
// Component specific styles
</style>
```

### Imports Organization
Order your `<script setup>` imports consistently:
1. Vue core (`ref`, `computed`, `onMounted`)
2. Vue Router (`useRouter`, `useRoute`)
3. Wails auto-generated bindings (`import * as ScriptService from '../../wailsjs/go/services/ScriptService'`)
4. Quasar and other UI libraries (`useQuasar`)
5. Local components (`../components/MyComponent.vue`)

### TypeScript & Types
- **Strict Mode** is enabled. Type everything explicitly.
- Define props using generics: `defineProps<{ title: string, count?: number }>()`.
- Define emits strongly: `defineEmits<{ (e: 'update', val: string): void }>()`.
- Wails generated models: `import { models } from '../../wailsjs/go/models'`. Use these for data passed from Go.

### UI & Styling
- Exclusively use **Quasar** components (`q-page`, `q-btn`, `q-dialog`, `q-input`).
- Dark mode overrides should use `body.body--dark` selectors and `:deep()`.
- SCSS is standard; use `<style lang="scss" scoped>`. Avoid inline CSS-in-JS and Tailwind.
- Use Quasar plugins for notifications: `$q.notify({ type: 'positive', message: 'Success' })`.
- UI text must be in **Chinese (Simplified)**. Keep variable names, function names, and comments in English.

### Async Data Fetching
Standardize API calls with try-catch-finally blocks:
```ts
const loadData = async () => {
  loading.value = true
  try {
    data.value = await SomeService.GetData()
  } catch (err) {
    console.error('Data load error:', err)
    $q.notify({ type: 'negative', message: '加载失败' })
  } finally {
    loading.value = false
  }
}
```

## Critical Rules & Pitfalls
1. **Never edit `frontend/wailsjs/` manually:** It is overwritten by Wails. Regenerate bindings after changing Go code by running `wails dev` or `wails generate module`.
2. **Platform prerequisites:** Linux builds require GTK/WebKit dev headers (`libgtk-3-dev`, `libwebkit2gtk-4.0-dev`).
3. **Database context:** Images are served locally via `MediaServer`; their paths are stored as JSON arrays in `Script.Images`.
