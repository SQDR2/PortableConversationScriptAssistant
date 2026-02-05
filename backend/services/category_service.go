package services

import (
	"context"
	"sidekick/backend/db"
	"sidekick/backend/models"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"gorm.io/gorm"
)

type CategoryService struct {
	ctx context.Context
}

func NewCategoryService() *CategoryService {
	return &CategoryService{}
}

func (s *CategoryService) Startup(ctx context.Context) {
	s.ctx = ctx
}

func (s *CategoryService) CreateCategory(name string) (*models.Category, error) {
	category := &models.Category{
		Name: name,
	}
	result := db.DB.Create(category)
	if result.Error != nil {
		runtime.LogErrorf(s.ctx, "Failed to create category: %v", result.Error)
		return nil, result.Error
	}
	return category, nil
}

func (s *CategoryService) ListCategories() ([]models.Category, error) {
	var categories []models.Category
	// Preload scripts count if needed, but for now just list categories
	// We might want to order by creation time or name
	result := db.DB.Order("created_at desc").Find(&categories)
	return categories, result.Error
}

func (s *CategoryService) UpdateCategory(id uint, name string) error {
	return db.DB.Model(&models.Category{}).Where("id = ?", id).Update("name", name).Error
}

func (s *CategoryService) DeleteCategory(id uint, cascade bool) error {
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
