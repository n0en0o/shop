package entities

type CatalogItem struct {
	BaseEntity
	ShortDescription string    `json:"short_description"`
	FullDescription  string    `json:"full_description"`
	ImageURL         string    `json:"image_url"`
	Brand            *Brand    `json:"brand"`
	Category         *Category `json:"category"`
	Price            float64   `json:"price"`
}
