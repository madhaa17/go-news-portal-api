package pagination

import (
	"math"

	"news-app/internal/core/domain/entity"
)

type PaginationInterface interface {
	AddPagination(totalData, page, perPage int) (*entity.Page, error)
}

type Options struct{}

// AddPagination implements PaginationInterface.
func (o *Options) AddPagination(totalData int, page int, perPage int) (*entity.Page, error) {
	newPage := page

	if newPage <= 0 {
		return nil, ErrorPage
	}

	limitData := 10

	if perPage > 0 {
		limitData = perPage
	}

	totalPage := math.Ceil(float64(totalData) / float64(limitData))

	last := (newPage * limitData)
	first := last - limitData

	if totalData < last {
		last = totalData
	}

	zeroPage := &entity.Page{PageCount: 1, Page: newPage}
	if totalData == 0 && page == 1 {
		return zeroPage, nil
	}

	if newPage > int(totalPage) {
		return nil, ErrorMaxPage
	}

	pages := &entity.Page{
		Page:       newPage,
		Perpage:    perPage,
		PageCount:  int(totalPage),
		TotalCount: totalData,
		First:      first,
		Last:       last,
	}

	return pages, nil
}

func NewPagination() PaginationInterface {
	pagination := new(Options)
	return pagination
}
