package pagination

import "errors"

var (
	ErrorMaxPage     = errors.New("choosen page more than total page")
	ErrorPage        = errors.New("page must be greater than 0")
	ErrorPageEmpty   = errors.New("page can't be empty")
	ErrorPageInvalid = errors.New("page must be a number")
)
