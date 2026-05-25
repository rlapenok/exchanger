package http

// PaginationQuery is the query parameters for the pagination
type PaginationQuery struct {
	Offset int `form:"offset"`
	Limit  int `form:"limit"`
}
