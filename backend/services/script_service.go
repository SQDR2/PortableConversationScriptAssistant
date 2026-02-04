package services

import (
	"context"
	"sidekick/backend/db"
	"sidekick/backend/models"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type ScriptService struct {
	ctx context.Context
}

func NewScriptService() *ScriptService {
	return &ScriptService{}
}

func (s *ScriptService) Startup(ctx context.Context) {
	s.ctx = ctx
	// Auto migrate
	if db.DB != nil {
		err := db.DB.AutoMigrate(&models.Script{})
		if err != nil {
			runtime.LogErrorf(s.ctx, "Failed to migrate database: %v", err)
		}
		s.setupFTS()
	}
}

func (s *ScriptService) setupFTS() {
	// Create virtual table
	db.DB.Exec("CREATE VIRTUAL TABLE IF NOT EXISTS scripts_fts USING fts5(content, tags, content='scripts', content_rowid='id');")

	// Create triggers to keep FTS index in sync
	db.DB.Exec(`CREATE TRIGGER IF NOT EXISTS scripts_ai AFTER INSERT ON scripts BEGIN
        INSERT INTO scripts_fts(rowid, content, tags) VALUES (new.id, new.content, new.tags);
    END;`)
	db.DB.Exec(`CREATE TRIGGER IF NOT EXISTS scripts_ad AFTER DELETE ON scripts BEGIN
        INSERT INTO scripts_fts(scripts_fts, rowid, content, tags) VALUES('delete', old.id, old.content, old.tags);
    END;`)
	db.DB.Exec(`CREATE TRIGGER IF NOT EXISTS scripts_au AFTER UPDATE ON scripts BEGIN
        INSERT INTO scripts_fts(scripts_fts, rowid, content, tags) VALUES('delete', old.id, old.content, old.tags);
        INSERT INTO scripts_fts(rowid, content, tags) VALUES (new.id, new.content, new.tags);
    END;`)
}

func (s *ScriptService) CreateScript(content string, tags string) (*models.Script, error) {
	script := &models.Script{
		Content: content,
		Tags:    tags,
	}
	result := db.DB.Create(script)
	return script, result.Error
}

func (s *ScriptService) UpdateScript(id uint, content string, tags string) (*models.Script, error) {
	var script models.Script
	if err := db.DB.First(&script, id).Error; err != nil {
		return nil, err
	}
	script.Content = content
	script.Tags = tags
	if err := db.DB.Save(&script).Error; err != nil {
		return nil, err
	}
	return &script, nil
}

func (s *ScriptService) DeleteScript(id uint) error {
	return db.DB.Delete(&models.Script{}, id).Error
}

func (s *ScriptService) ListScripts(page int, pageSize int) ([]models.Script, error) {
	var scripts []models.Script
	offset := (page - 1) * pageSize
	result := db.DB.Order("created_at desc").Offset(offset).Limit(pageSize).Find(&scripts)
	return scripts, result.Error
}

func (s *ScriptService) SearchScripts(query string) ([]models.Script, error) {
	var scripts []models.Script
	// Use FTS5 match query
	// Ideally we sanitize query or use separate args, strictly FTS syntax
	// GORM raw query
	// Note: We select from the MAIN table, joined with FTS, or just use rank
	// Simplest: SELECT * FROM scripts WHERE id IN (SELECT rowid FROM scripts_fts WHERE scripts_fts MATCH ?)

	err := db.DB.Raw("SELECT * FROM scripts WHERE id IN (SELECT rowid FROM scripts_fts WHERE scripts_fts MATCH ?) ORDER BY rank", query).Scan(&scripts).Error
	if err != nil {
		runtime.LogErrorf(s.ctx, "Search error: %v", err)
		return nil, err
	}
	return scripts, nil
}

func (s *ScriptService) ImportScripts(content string, delimiter string) (int, error) {
	parts := strings.Split(content, delimiter)
	count := 0
	tx := db.DB.Begin()
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed == "" {
			continue
		}
		script := &models.Script{Content: trimmed}
		if err := tx.Create(script).Error; err != nil {
			tx.Rollback()
			return 0, err
		}
		count++
	}
	return count, tx.Commit().Error
}
