package dto

type IDResponse[T int64 | int32 | string] struct {
	ID T `json:"id"`
}

// PagingResponse is a generic paged response
type PagingResponse[T any] struct {
	Items      []T   `json:"items"`       // paged items
	Total      int64 `json:"total"`       // total number of items
	Page       int   `json:"page"`        // current page
	PageSize   int   `json:"page_size"`   // number of items per page
	TotalPages int   `json:"total_pages"` // total pages
}
