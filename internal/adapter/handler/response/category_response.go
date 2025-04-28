package response

type SuccessCategoryResponse struct {
	ID            int16  `json:"id"`
	Title         string `json:"title"`
	Slug          string `json:"slug"`
	CreatedByName string `json:"created_by_name"`
}
