package repository

import (
	"context"
	"errors"
	"fmt"

	"news-app/internal/core/domain/entity"
	"news-app/internal/core/domain/model"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	GetCategories(ctx context.Context) ([]entity.CategoryEntity, error)
	GetCategoryByID(ctx context.Context, id int16) (*entity.CategoryEntity, error)
	CreateCategory(ctx context.Context, req entity.CategoryEntity) (*entity.CategoryEntity, error)
	UpdateCategory(ctx context.Context, req entity.CategoryEntity) (*entity.CategoryEntity, error)
	DeleteCategory(ctx context.Context, id int16) error
}

type categoryRepository struct {
	db *gorm.DB
}

// CreateCategory implements CategoryRepository.
func (c *categoryRepository) CreateCategory(ctx context.Context, req entity.CategoryEntity) (*entity.CategoryEntity, error) {
	var countSlug int64
	err = c.db.Table("categories").Where("slug = ?", req.Slug).Count(&countSlug).Error
	if err != nil {
		code := "[REPOSITORY] CreateCategory - 1"
		log.Errorw(code, err)
		return nil, err
	}

	countSlug = countSlug + 1
	slug := fmt.Sprintf("%s-%d", req.Slug, countSlug)
	modelCategory := model.Category{
		Title:       req.Title,
		Slug:        slug,
		CreatedByID: int64(req.UserEntity.ID),
	}

	err = c.db.Create(&modelCategory).Error
	if err != nil {
		code = "[REPOSITORY] CreateCategory - 2"
		log.Errorw(code, err)
		return nil, err
	}

	return nil, err
}

// DeleteCategory implements CategoryRepository.
func (c *categoryRepository) DeleteCategory(ctx context.Context, id int16) error {
	var count int64

	err = c.db.Table("contents").Where("category_id = ?", id).Count(&count).Error
	if err != nil {
		code := "[REPOSITORY] DeleteCategory - 1"
		log.Errorw(code, err)
		return err
	}

	if count > 0 {
		return errors.New("cannot delete a category that has  associated contents")
	}

	err = c.db.Where("id = ?", id).Delete(&model.Category{}).Error
	if err != nil {
		code := "[REPOSITORY] DeleteCategory - 2"
		log.Errorw(code, err)
		return err
	}

	return nil
}

// GetCategories implements CategoryRepository.
func (c *categoryRepository) GetCategories(ctx context.Context) ([]entity.CategoryEntity, error) {
	var modelCategories []model.Category

	err := c.db.Order("created_at DESC").Preload("User").Find(&modelCategories).Error
	if err != nil {
		code = "[REPOSITORY] GetCategories - 1"
		log.Errorw(code, err)
		return nil, err
	}

	if len(modelCategories) == 0 {
		code = "[REPOSITORY] GetCategories - 2"
		err = errors.New("data not found")
		log.Errorw(code, err)
		return nil, err
	}

	var res []entity.CategoryEntity
	for _, val := range modelCategories {
		res = append(res, entity.CategoryEntity{
			ID:    int16(val.ID),
			Title: val.Title,
			Slug:  val.Slug,
			UserEntity: entity.UserEntity{
				ID:       int16(val.User.ID),
				Name:     val.User.Name,
				Email:    val.User.Email,
				Password: val.User.Password,
			},
		})
	}

	return res, nil
}

// GetCategoryByID implements CategoryRepository.
func (c *categoryRepository) GetCategoryByID(ctx context.Context, id int16) (*entity.CategoryEntity, error) {
	var modelCategory model.Category

	err := c.db.Where("id = ?", id).Preload("User").First(&modelCategory).Error
	if err != nil {
		code := "[REPOSITORY] GetCategoryByID - 1"
		log.Errorw(code, err)
		return nil, err
	}

	return &entity.CategoryEntity{
		ID:    int16(modelCategory.ID),
		Title: modelCategory.Title,
		Slug:  modelCategory.Slug,
		UserEntity: entity.UserEntity{
			ID:    int16(modelCategory.User.ID),
			Name:  modelCategory.User.Name,
			Email: modelCategory.User.Email,
		},
	}, err
}

// UpdateCategory implements CategoryRepository.
func (c *categoryRepository) UpdateCategory(ctx context.Context, req entity.CategoryEntity) (*entity.CategoryEntity, error) {
	var countSlug int64
	err := c.db.Table("categories").Where("slug = ?", req.Slug).Count(&countSlug).Error
	if err != nil {
		code := "[REPOSITORY] UpdateCategory - 1"
		log.Errorw(code, err)
		return nil, err
	}

	countSlug = countSlug + 1
	slug := req.Slug

	if countSlug == 1 {
		slug = fmt.Sprintf("%s-%d", req.Slug, countSlug)
	}

	modelCategory := model.Category{
		Title:       req.Title,
		Slug:        slug,
		CreatedByID: int64(req.UserEntity.ID),
	}

	err = c.db.Where("id = ?", req.ID).Updates(&modelCategory).Error
	if err != nil {
		code := "[REPOSITORY] UpdateCategory - 2"
		log.Errorw(code, err)
		return nil, err
	}

	return nil, err
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}
