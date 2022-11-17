package user

import (
	"backend/businesses/users"
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint				`json:"id" gorm:"primaryKey"`
	CreatedAt time.Time			`json:"created_at"`
	UpdatedAt time.Time			`json:"updated_at"`
	DeletedAt gorm.DeletedAt	`json:"deleted_at"`
	Fullname  string			`json:"fullname"`
	Email     string		 	`json:"email" gorm:"unique" faker:"email"`
	Password  string			`json:"password" faker:"password"`
	Gender    string			`json:"gender"`
	DOB       time.Time			`json:"dob"`
	Phone     string			`json:"phone"`
}

func FromDomain(domain *user.Domain) *User {
	return &User{
		ID: domain.ID,
		CreatedAt: domain.CreatedAt,
		UpdatedAt: domain.UpdatedAt,
		DeletedAt: domain.DeletedAt,
		Fullname: domain.Fullname,
		Email: domain.Email,
		Password: domain.Email,
		Gender: domain.Gender,
		DOB: domain.DOB,
		Phone: domain.Phone,
	}
}

func (rec *User) ToDomain() user.Domain {
	return user.Domain{
		ID: rec.ID,
		Fullname: rec.Fullname,
		Email: rec.Email,
		Password: rec.Password,
		Gender: rec.Gender,
		DOB: rec.DOB,
		Phone: rec.Phone,
		CreatedAt: rec.CreatedAt,
		UpdatedAt: rec.UpdatedAt,
		DeletedAt: rec.DeletedAt,
	}
}