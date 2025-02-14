package domain

type Customer struct {
	ID            uint   `json:"id" gorm:"primaryKey"`
	Name          string `json:"name"`
	Email         string `json:"email" gorm:"unique"`
	VerifiedEmail bool   `json:"verified_email"`
	Phone         string `json:"phone"`
	Picture       string `json:"picture"`
}
