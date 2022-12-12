package transactions

import (
	"time"

	"gorm.io/gorm"
)

type Domain struct {
	ID            uint
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt
	Price         uint
	CheckIn       time.Time
	CheckOut      time.Time
	Duration      int
	PaymentMethod string
	Status        string
	Drink         string
	UserID        uint
	UserFullName  string
	UserEmail     string
	OfficeID      uint
	OfficeName    string
}

type Usecase interface {
	GetAll() []Domain
	GetByUserID(userId string) []Domain
	AdminGetByUserID(userId string) []Domain
	GetByOfficeID(officeId string) []Domain
	Create(transactionDomain *Domain) Domain
	GetByID(id string) Domain
	Update(id string, transactionDomain *Domain) Domain
	Delete(id string) bool
}

type Repository interface {
	GetAll() []Domain
	GetByUserID(userId string) []Domain
	GetByOfficeID(officeId string) []Domain
	Create(transactionDomain *Domain) Domain
	GetByID(id string) Domain
	Update(id string, transactionDomain *Domain) Domain
	Delete(id string) bool
}
