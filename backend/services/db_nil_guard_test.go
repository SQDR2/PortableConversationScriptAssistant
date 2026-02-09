package services

import (
	"testing"

	"sidekick/backend/db"
)

func TestListScripts_DBNil_ReturnsEmptySliceAndError_NoPanic(t *testing.T) {
	prevDB := db.DB
	defer func() { db.DB = prevDB }()
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("unexpected panic: %v", r)
		}
	}()

	db.DB = nil
	svc := NewScriptService() // ctx intentionally nil
	scripts, err := svc.ListScripts(1, 20)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if scripts == nil {
		t.Fatalf("expected empty slice, got nil")
	}
	if len(scripts) != 0 {
		t.Fatalf("expected 0 scripts, got %d", len(scripts))
	}
}

func TestListCategories_DBNil_ReturnsEmptySliceAndError_NoPanic(t *testing.T) {
	prevDB := db.DB
	defer func() { db.DB = prevDB }()
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("unexpected panic: %v", r)
		}
	}()

	db.DB = nil
	svc := NewCategoryService() // ctx intentionally nil
	categories, err := svc.ListCategories()
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if categories == nil {
		t.Fatalf("expected empty slice, got nil")
	}
	if len(categories) != 0 {
		t.Fatalf("expected 0 categories, got %d", len(categories))
	}
}
