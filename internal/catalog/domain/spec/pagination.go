package spec

type Pagination[T any] struct {
	PageIndex  int `json:"page_index"`
	PageSize   int `json:"page_size"`
	TotalCount int `json:"total_count"`
	Items      []T `json:"items"`
}
