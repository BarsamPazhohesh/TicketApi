package dto

type IDResponse[T int64 | int32 | string] struct {
	ID T `json:"id"`
}
