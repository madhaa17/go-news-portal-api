package service

import (
	"context"

	"news-app/internal/adapter/repository"
	"news-app/internal/core/domain/entity"
	"news-app/lib/conv"

	"github.com/gofiber/fiber/v2/log"
)

type CategoryService interface {
	GetCategories(ctx context.Context) ([]entity.CategoryEntity, error)
	GetCategoryByID(ctx context.Context, id int16) (*entity.CategoryEntity, error)
	CreateCategory(ctx context.Context, req entity.CategoryEntity) (*entity.CategoryEntity, error)
	UpdateCategory(ctx context.Context, req entity.CategoryEntity) (*entity.CategoryEntity, error)
	DeleteCategory(ctx context.Context, id int16) error
}

type categoryService struct {
	categoryRepository repository.CategoryRepository
}

// CreateCategory implements CategoryService.
func (c *categoryService) CreateCategory(ctx context.Context, req entity.CategoryEntity) (*entity.CategoryEntity, error) {
	slug := conv.GeneratesSlug(req.Title)
	req.Slug = slug

	_, err := c.categoryRepository.CreateCategory(ctx, req)
	if err != nil {
		code := "[SERVICE] CreateCategory - 1"
		log.Errorw(code, err)
		return nil, err
	}

	return &req, nil
}

// DeleteCategory implements CategoryService.
func (c *categoryService) DeleteCategory(ctx context.Context, id int16) error {
	err = c.categoryRepository.DeleteCategory(ctx, id)
	if err != nil {
		code := "[SERVICE] DeleteCategory - 1"
		log.Errorw(code, err)
		return err
	}

	return nil
}

// GetCategories implements CategoryService.
func (c *categoryService) GetCategories(ctx context.Context) ([]entity.CategoryEntity, error) {
	results, err := c.categoryRepository.GetCategories(ctx)
	if err != nil {
		code = "[SERVICE] GetCategories - 1"
		log.Errorw(code, err)
		return nil, err
	}

	return results, nil
}

// GetCategoryByID implements CategoryService.
func (c *categoryService) GetCategoryByID(ctx context.Context, id int16) (*entity.CategoryEntity, error) {
	results, err := c.categoryRepository.GetCategoryByID(ctx, id)
	if err != nil {
		code = "[SERVICE] GetCategoryByID - 1"
		log.Errorw(code, err)
		return nil, err
	}

	return results, nil
}

// UpdateCategory implements CategoryService.
func (c *categoryService) UpdateCategory(ctx context.Context, req entity.CategoryEntity) (*entity.CategoryEntity, error) {
	categoryData, err := c.categoryRepository.GetCategoryByID(ctx, req.ID)
	if err != nil {
		code := "[SERVICE] UpdateCategory - 1"
		log.Errorw(code, err)
		return nil, err
	}
	slug := conv.GeneratesSlug(req.Title)

	if categoryData.Title == req.Title {
		req.Slug = categoryData.Slug
	}

	req.Slug = slug

	_, err = c.categoryRepository.UpdateCategory(ctx, req)
	if err != nil {
		code := "[SERVICE] UpdateCategory - 2"
		log.Errorw(code, err)
		return nil, err
	}

	return nil, err
}

func NewCategoryService(categoryRepo repository.CategoryRepository) CategoryService {
	return &categoryService{categoryRepository: categoryRepo}
}
