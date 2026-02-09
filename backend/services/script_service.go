package services

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sidekick/backend/db"
	"sidekick/backend/models"
	"strings"

	"github.com/google/uuid"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type ScriptService struct {
	ctx context.Context
}

func (s *ScriptService) logInfof(format string, args ...interface{}) {
	if s.ctx != nil {
		runtime.LogInfof(s.ctx, format, args...)
		return
	}
	log.Printf("[ScriptService] "+format, args...)
}

func (s *ScriptService) logInfo(msg string) {
	if s.ctx != nil {
		runtime.LogInfo(s.ctx, msg)
		return
	}
	log.Printf("[ScriptService] %s", msg)
}

func (s *ScriptService) logErrorf(format string, args ...interface{}) {
	if s.ctx != nil {
		runtime.LogErrorf(s.ctx, format, args...)
		return
	}
	log.Printf("[ScriptService][ERROR] "+format, args...)
}

func NewScriptService() *ScriptService {
	return &ScriptService{}
}

func (s *ScriptService) Startup(ctx context.Context) {
	s.ctx = ctx
	s.logInfo("ScriptService starting up...")
	// Auto migrate
	if db.DB != nil {
		err := db.DB.AutoMigrate(&models.Category{}, &models.Script{})
		if err != nil {
			s.logErrorf("Failed to migrate database: %v", err)
		}
		s.logInfo("Database migration complete")
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

func (s *ScriptService) CreateScript(content string, tags string, categoryID *uint, images string) (*models.Script, error) {
	if db.DB == nil {
		err := fmt.Errorf("database not initialized")
		s.logErrorf("CreateScript failed: %v", err)
		return nil, err
	}
	script := &models.Script{
		Content:    content,
		Tags:       tags,
		CategoryID: categoryID,
		Images:     images,
	}
	result := db.DB.Create(script)
	return script, result.Error
}

func (s *ScriptService) UpdateScript(id uint, content string, tags string, categoryID *uint, images string) (*models.Script, error) {
	if db.DB == nil {
		err := fmt.Errorf("database not initialized")
		s.logErrorf("UpdateScript failed: %v", err)
		return nil, err
	}
	var script models.Script
	if err := db.DB.First(&script, id).Error; err != nil {
		return nil, err
	}
	script.Content = content
	script.Tags = tags
	script.CategoryID = categoryID
	script.Images = images
	if err := db.DB.Save(&script).Error; err != nil {
		return nil, err
	}
	return &script, nil
}

func (s *ScriptService) DeleteScript(id uint) error {
	if db.DB == nil {
		err := fmt.Errorf("database not initialized")
		s.logErrorf("DeleteScript failed: %v", err)
		return err
	}
	var script models.Script
	if err := db.DB.First(&script, id).Error; err == nil {
		if script.Images != "" {
			var imagePaths []string
			if err := json.Unmarshal([]byte(script.Images), &imagePaths); err == nil {
				for _, p := range imagePaths {
					cwd, _ := os.Getwd()
					fullPath := filepath.Join(cwd, p)
					os.Remove(fullPath)
				}
			}
		}
	}
	return db.DB.Delete(&models.Script{}, id).Error
}

func (s *ScriptService) ListScripts(page int, pageSize int) ([]models.Script, error) {
	s.logInfof("ListScripts called: page=%d, pageSize=%d", page, pageSize)
	scripts := make([]models.Script, 0)
	if db.DB == nil {
		err := fmt.Errorf("database not initialized")
		s.logErrorf("ListScripts failed: %v", err)
		return scripts, err
	}
	offset := (page - 1) * pageSize
	result := db.DB.Order("created_at desc").Offset(offset).Limit(pageSize).Find(&scripts)
	if result.Error != nil {
		s.logErrorf("ListScripts DB error: %v", result.Error)
	}
	s.logInfof("ListScripts found %d records", len(scripts))
	return scripts, result.Error
}

func (s *ScriptService) SearchScripts(query string) ([]models.Script, error) {
	scripts := make([]models.Script, 0)
	if db.DB == nil {
		return scripts, fmt.Errorf("database not initialized")
	}
	searchQuery := "%" + query + "%"
	// Standard LIKE search
	result := db.DB.Where("content LIKE ? OR tags LIKE ?", searchQuery, searchQuery).
		Order("created_at desc").
		Find(&scripts)

	return scripts, result.Error
}

func (s *ScriptService) ImportScripts(scripts []string) (int, error) {
	if db.DB == nil {
		return 0, fmt.Errorf("database not initialized")
	}
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

func (s *ScriptService) SaveScriptImage(base64Data string, ext string) (string, error) {
	// Decode base64
	data, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %v", err)
	}

	// Ensure images directory exists
	cwd, _ := os.Getwd()
	imagesDir := filepath.Join(cwd, "images")
	if err := os.MkdirAll(imagesDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create images directory: %v", err)
	}

	// Generate UUID V7 filename
	id, err := uuid.NewV7()
	if err != nil {
		return "", fmt.Errorf("failed to generate uuid: %v", err)
	}
	fileName := id.String() + ext
	filePath := filepath.Join(imagesDir, fileName)

	// Save file
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return "", fmt.Errorf("failed to write file: %v", err)
	}

	// Return RELATIVE path for database
	return filepath.Join("images", fileName), nil
}
