package spec

import "github.com/google/uuid"

const MaxPageSize = 4

type QueryArgs struct {
	PageIndex  int     `form:"page_index"`
	PageSize   int     `form:"page_size"`
	BrandID    *string `form:"brand_id"`
	CategoryID *string `form:"category_id"`
	Search     *string `form:"search"`
	Sort       *string `form:"sort"`
}

func (q *QueryArgs) Normalize() {

	if q.PageIndex <= 0 {
		q.PageIndex = 1
	}

	if q.PageSize <= 0 {
		q.PageSize = 2
	}

	if q.PageSize > MaxPageSize {
		q.PageSize = MaxPageSize
	}

}

func (q *QueryArgs) ParseBrandID() (*uuid.UUID, error) {
	if q.BrandID == nil || *q.BrandID == "" {
		return nil, nil
	}

	id, err := uuid.Parse(*q.BrandID)
	if err != nil {
		return nil, err
	}
	return &id, nil
}

func (q *QueryArgs) ParseCategoryID() (*uuid.UUID, error) {
	if q.CategoryID == nil || *q.CategoryID == "" {
		return nil, nil
	}

	id, err := uuid.Parse(*q.CategoryID)
	if err != nil {
		return nil, err
	}
	return &id, nil
}
