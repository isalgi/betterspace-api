package response

import (
	"backend/businesses/users"
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
	FullName  string         `json:"full_name"`
	Gender    string         `json:"gender"`
	Email     string         `json:"email"`
	Photo     string         `json:"photo" form:"photo"`
}

func FromDomain(domain users.Domain) User {
	return User{
		ID:        domain.ID,
		CreatedAt: domain.CreatedAt,
		UpdatedAt: domain.UpdatedAt,
		DeletedAt: domain.DeletedAt,
		FullName:  domain.FullName,
		Gender:    domain.Gender,
		Email:     domain.Email,
		Photo:     domain.Photo,
	}
}