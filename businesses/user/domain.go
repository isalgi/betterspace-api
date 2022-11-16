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
	DOB       string
	Phone     int
}

type Usecase interface {
	Login(userDomain *Domain) string
	Register(userDomain *Domain) Domain
}

type Repository interface {
	Register(userDomain *Domain) Domain
}
