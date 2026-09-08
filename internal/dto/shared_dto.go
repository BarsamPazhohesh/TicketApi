package dto

type IDResponse[T int64 | int32 | string] struct {
	ID T `json:"id"`
}

type IDResponseInt64 struct {
	ID int64 `json:"id"`
}

type IDResponseString struct {
	ID string `json:"id"`
}

type IDRequest[T int64 | int32 | string] = IDResponse[T]
type IDRequestInt64 = IDResponseInt64
type IDRequestString = IDResponseString

// PagingResponse is a generic paged response
type PagingResponse[T any] struct {
	Items      []T   `json:"items"`       // paged items
	Total      int64 `json:"total"`       // total number of items
	Page       int   `json:"page"`        // current page
	PageSize   int   `json:"page_size"`   // number of items per page
	TotalPages int   `json:"total_pages"` // total pages
}
