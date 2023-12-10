package models

type ContentPagination[T any] struct {
	Content []T
	Count   int64
}
