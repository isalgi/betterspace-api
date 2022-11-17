package user

import (
	"time"

	"gorm.io/gorm"
)

type Domain struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
	Fullname  string
	Email     string
	Password  string
	Gender    string
	DOB       time.Time
	Phone     string
}

type Usecase interface {
	Register(userDomain *Domain) Domain
	Login(userDomain *Domain) string
}

type Repository interface {
	Register(userDomain *Domain) Domain
	GetByEmail(userDomain *Domain) Domain
}
