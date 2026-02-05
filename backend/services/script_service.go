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
		err := db.DB.AutoMigrate(&models.Category{}, &models.Script{})
		if err != nil {
			runtime.LogErrorf(s.ctx, "Failed to migrate database: %v", err)
		}
		s.setupFTS()
	}
}

func (s *ScriptService) setupFTS() {
	// FTS5 might not be available on all systems, fallback to standard LIKE search.
	// Clean up any existing FTS5 artifacts to prevent errors.
	db.DB.Exec("DROP TRIGGER IF EXISTS scripts_ai")
	db.DB.Exec("DROP TRIGGER IF EXISTS scripts_ad")
	db.DB.Exec("DROP TRIGGER IF EXISTS scripts_au")
	db.DB.Exec("DROP TABLE IF EXISTS scripts_fts")
}

func (s *ScriptService) CreateScript(content string, tags string, categoryID *uint) (*models.Script, error) {
	script := &models.Script{
		Content:    content,
		Tags:       tags,
		CategoryID: categoryID,
	}
	result := db.DB.Create(script)
	return script, result.Error
}

func (s *ScriptService) UpdateScript(id uint, content string, tags string, categoryID *uint) (*models.Script, error) {
	var script models.Script
	if err := db.DB.First(&script, id).Error; err != nil {
		return nil, err
	}
	script.Content = content
	script.Tags = tags
	script.CategoryID = categoryID
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
	searchQuery := "%" + query + "%"
	// Standard LIKE search
	result := db.DB.Where("content LIKE ? OR tags LIKE ?", searchQuery, searchQuery).
		Order("created_at desc").
		Find(&scripts)

	return scripts, result.Error
}

func (s *ScriptService) ImportScripts(scripts []string) (int, error) {
	count := 0
	tx := db.DB.Begin()
	for _, part := range scripts {
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
