package domain

type Product struct {
	ID         uint     `json:"id"`
	Name       string   `json:"name"`
	Price      float64  `json:"price"`
	CategoryID uint     `json:"category_id"`
	Category   Category `json:"category" gormm:"foreignKey:CategoryID"`
}
