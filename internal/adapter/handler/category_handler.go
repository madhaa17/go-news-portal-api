package handler

import (
	"news-app/internal/adapter/handler/request"
	"news-app/internal/adapter/handler/response"
	"news-app/internal/core/domain/entity"
	"news-app/internal/core/service"
	validatorLib "news-app/lib/validator"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

var defaultResponse response.DefaultSuccessResponse

type CategoryHandler interface {
	GetCategories(c *fiber.Ctx) error
	GetCategoryByID(c *fiber.Ctx) error
	CreateCategory(c *fiber.Ctx) error
	UpdateCategory(c *fiber.Ctx) error
	DeleteCategory(c *fiber.Ctx) error
}

type categoryHandler struct {
	categoryService service.CategoryService
}

// CreateCategory implements CategoryHandler.
func (ch *categoryHandler) CreateCategory(c *fiber.Ctx) error {
	var req request.CategoryRequest
	claims := c.Locals("user").(*entity.JwtData)
	userID := claims.UserID

	if userID == 0 {
		code = "[HANDLER] CreateCategory - 1"
		log.Errorw(code, err)
		errResponse.Meta.Status = false
		errResponse.Meta.Message = "Unauthorized access"

		return c.Status(fiber.StatusUnauthorized).JSON(errResponse)
	}

	if err := c.BodyParser(&req); err != nil {
		code = "[HANDLER] CreateCategory - 2"
		log.Errorw(code, err)
		errResponse.Meta.Status = false
		errResponse.Meta.Message = err.Error()

		return c.Status(fiber.StatusBadRequest).JSON(errResponse)
	}

	if err := validatorLib.ValidateStruct(req); err != nil {
		code = "[HANDLER] CreateCategory - 3"
		log.Errorw(code, err)
		errResponse.Meta.Status = false
		errResponse.Meta.Message = err.Error()

		return c.Status(fiber.StatusBadRequest).JSON(errResponse)
	}

	reqEntity := entity.CategoryEntity{
		Title: req.Title,
		UserEntity: entity.UserEntity{
			ID: int16(userID),
		},
	}

	_, err = ch.categoryService.CreateCategory(c.Context(), reqEntity)
	if err != nil {
		code = "[HANDLER] CreateCategory - 4"
		log.Errorw(code, err)
		errResponse.Meta.Status = false
		errResponse.Meta.Message = err.Error()

		return c.JSON(defaultResponse)
	}

	defaultResponse.Data = nil
	defaultResponse.Pagination = nil
	defaultResponse.Meta.Status = true
	defaultResponse.Meta.Message = "Category created successfully"
	return c.JSON(defaultResponse)
}

// DeleteCategory implements CategoryHandler.
func (ch *categoryHandler) DeleteCategory(c *fiber.Ctx) error {
	panic("unimplemented")
}

// GetCategories implements CategoryHandler.
func (ch *categoryHandler) GetCategories(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	userID := claims.UserID

	log.Infof("User ID: %s", userID)

	if userID == 0 {
		code = "[HANDLER] GetCategories - 1"
		log.Errorw(code, err)
		errResponse.Meta.Status = false
		errResponse.Meta.Message = "Unauthorized access"

		return c.Status(fiber.StatusUnauthorized).JSON(errResponse)
	}

	results, err := ch.categoryService.GetCategories(c.Context())
	if err != nil {
		code = "[HANDLER] GetCategories - 2"
		log.Errorw(code, err)
		errResponse.Meta.Status = false
		errResponse.Meta.Message = err.Error()

		return c.Status(fiber.StatusInternalServerError).JSON(errResponse)
	}

	categoryReponses := []response.SuccessCategoryResponse{}
	for _, result := range results {
		categoryResponse := response.SuccessCategoryResponse{
			ID:            result.ID,
			Title:         result.Title,
			Slug:          result.Slug,
			CreatedByName: result.UserEntity.Name,
		}

		categoryReponses = append(categoryReponses, categoryResponse)
	}

	defaultResponse.Meta.Status = true
	defaultResponse.Pagination = nil
	defaultResponse.Meta.Message = "Categories fetched successfully"
	defaultResponse.Data = categoryReponses

	return c.JSON(defaultResponse)
}

// GetCategoryByID implements CategoryHandler.
func (ch *categoryHandler) GetCategoryByID(c *fiber.Ctx) error {
	panic("unimplemented")
}

// UpdateCategory implements CategoryHandler.
func (ch *categoryHandler) UpdateCategory(c *fiber.Ctx) error {
	panic("unimplemented")
}

func NewCategoryHandler(categoryService service.CategoryService) CategoryHandler {
	return &categoryHandler{categoryService: categoryService}
}
