package uc

const (
	DefaultLimit  = 10
	MaxLimit      = 100
	DefaultOffset = 0
)

// Pagination is the pagination for the use case
type Pagination struct {
	Limit  int
	Offset int
}

// NewPagination creates a new Pagination
func NewPagination(limit, offset int) Pagination {
	if limit < 1 {
		limit = DefaultLimit
	}
	if limit > MaxLimit {
		limit = MaxLimit
	}
	if offset < 0 {
		offset = 0
	}
	return Pagination{
		Limit:  limit,
		Offset: offset,
	}
}
