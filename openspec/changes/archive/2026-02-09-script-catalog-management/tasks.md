## 1. Backend: Data Model & Service

- [x] 1.1 Create `Category` model in `models/category.go`
- [x] 1.2 Update `Script` model in `models/script.go` to include `CategoryID`
- [x] 1.3 Ensure `db.AutoMigrate` handles the schema changes in `script_service.go` or `main.go`
- [x] 1.4 Create `CategoryService` (in `services/category_service.go`)
  - [x] 1.4.1 Implement `CreateCategory`
  - [x] 1.4.2 Implement `ListCategories`
  - [x] 1.4.3 Implement `UpdateCategory`
  - [x] 1.4.4 Implement `DeleteCategory` (with cascade support)
- [x] 1.5 Update `ScriptService` to support categories
  - [x] 1.5.1 Update `CreateScript` to accept `categoryId`
  - [x] 1.5.2 Update `UpdateScript` to accept `categoryId`
- [x] 1.6 Register `CategoryService` in `main.go` for Wails binding

## 2. Frontend: Infrastructure & Components

- [x] 2.1 Regenerate Wails JS bindings (`wails dev` or manual trigger)
- [x] 2.2 Create `CategoryList` component
  - [x] Display list of categories
  - [x] Support "New Category" and "Edit Category" actions
- [x] 2.3 Create `ScriptItem` component
  - [x] Display content preview (no title)
  - [x] Implement "Copy" button
  - [x] Implement "Delete" button with confirmation dialog
  - [x] Implement "Edit" trigger
- [x] 2.4 Update `SelectTargetDialog` or `App.vue` (if necessary) to ensuring consistent styles (optional)

## 3. Frontend: Page Refactor (ScriptsPage.vue)

- [x] 3.1 Implement Tabs for View Switching (Timeline vs Directory)
- [x] 3.2 Implement Timeline View (using `ScriptItem` component)
- [x] 3.3 Implement Directory View
  - [x] Show categories as cards/list
  - [x] Show "Uncategorized" section
  - [x] Click category to filter script list
- [x] 3.4 Update "New/Edit Script" Dialog
  - [x] Add Category selection dropdown
- [x] 3.5 Implement Directory Management UI
  - [x] "New Directory" button and dialog
  - [x] Directory deletion with cascade confirmation dialog

## 4. Verification

- [ ] 4.1 Verify Timeline View displays all scripts by time
- [ ] 4.2 Verify creating, renaming, deleting categories
- [ ] 4.3 Verify deleting category with "Cascade" vs "Keep Scripts"
- [ ] 4.4 Verify creating/moving scripts into categories
- [ ] 4.5 Verify Copy/Edit/Delete script actions
