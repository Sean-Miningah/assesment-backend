package domain

type Category struct {
	ID               uint       `json:"id"`
	Name             string     `json:"name"`
	ParentCategoryID *uint      `json:"parent_category_id"`
	Category         []Category `json:"category" gorm:"foreignKey:ParentCategoryID"`
}
