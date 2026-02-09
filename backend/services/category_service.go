package services

import (
	"context"
	"fmt"
	"log"
	"sidekick/backend/db"
	"sidekick/backend/models"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"gorm.io/gorm"
)

type CategoryService struct {
	ctx context.Context
}

func (s *CategoryService) logInfo(msg string) {
	if s.ctx != nil {
		runtime.LogInfo(s.ctx, msg)
		return
	}
	log.Printf("[CategoryService] %s", msg)
}

func (s *CategoryService) logInfof(format string, args ...interface{}) {
	if s.ctx != nil {
		runtime.LogInfof(s.ctx, format, args...)
		return
	}
	log.Printf("[CategoryService] "+format, args...)
}

func (s *CategoryService) logErrorf(format string, args ...interface{}) {
	if s.ctx != nil {
		runtime.LogErrorf(s.ctx, format, args...)
		return
	}
	log.Printf("[CategoryService][ERROR] "+format, args...)
}

func NewCategoryService() *CategoryService {
	return &CategoryService{}
}

func (s *CategoryService) Startup(ctx context.Context) {
	s.ctx = ctx
}

func (s *CategoryService) CreateCategory(name string) (*models.Category, error) {
	if db.DB == nil {
		err := fmt.Errorf("database not initialized")
		s.logErrorf("CreateCategory failed: %v", err)
		return nil, err
	}
	category := &models.Category{
		Name: name,
	}
	result := db.DB.Create(category)
	if result.Error != nil {
		s.logErrorf("Failed to create category: %v", result.Error)
		return nil, result.Error
	}
	return category, nil
}

func (s *CategoryService) ListCategories() ([]models.Category, error) {
	s.logInfo("ListCategories called")
	categories := make([]models.Category, 0)
	if db.DB == nil {
		err := fmt.Errorf("database not initialized")
		s.logErrorf("ListCategories failed: %v", err)
		return categories, err
	}
	result := db.DB.Order("created_at desc").Find(&categories)
	if result.Error != nil {
		s.logErrorf("ListCategories DB error: %v", result.Error)
	}
	s.logInfof("ListCategories found %d records", len(categories))
	return categories, result.Error
}

func (s *CategoryService) UpdateCategory(id uint, name string) error {
	if db.DB == nil {
		err := fmt.Errorf("database not initialized")
		s.logErrorf("UpdateCategory failed: %v", err)
		return err
	}
	return db.DB.Model(&models.Category{}).Where("id = ?", id).Update("name", name).Error
}

func (s *CategoryService) DeleteCategory(id uint, cascade bool) error {
	if db.DB == nil {
		err := fmt.Errorf("database not initialized")
		s.logErrorf("DeleteCategory failed: %v", err)
		return err
	}
	return db.DB.Transaction(func(tx *gorm.DB) error {
		if cascade {
			// Delete scripts in this category
			if err := tx.Where("category_id = ?", id).Delete(&models.Script{}).Error; err != nil {
				return err
			}
		} else {
			// Move scripts to uncategorized (NULL)
			if err := tx.Model(&models.Script{}).Where("category_id = ?", id).Update("category_id", nil).Error; err != nil {
				return err
			}
		}
		// Delete the category itself
		if err := tx.Delete(&models.Category{}, id).Error; err != nil {
			return err
		}
		return nil
	})
}
