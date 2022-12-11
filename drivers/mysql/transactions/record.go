package transactions

import (
	TransactionUseCase "backend/businesses/transactions"
	"backend/drivers/mysql/offices"
	"backend/drivers/mysql/users"
	"time"

	"gorm.io/gorm"
)

type Transaction struct {
	ID       uint           `json:"id" gorm:"primaryKey"`
	CreatedAt	time.Time		`json:"created_at"`
	UpdatedAt	time.Time		`json:"updated_at"`
	DeletedAt	gorm.DeletedAt	`json:"deleted_at"`
	Price    uint           `json:"price"`
	UserID   uint           `json:"user_id"`
	OfficeID uint           `json:"office_id"`
	User     users.User     `json:"user" gorm:"foreignKey:UserID;references:ID"`
	Office   offices.Office `json:"office" gorm:"foreignKey:OfficeID;references:ID"`
}

func FromDomain(domain *TransactionUseCase.Domain) *Transaction {
	return &Transaction{
		ID:       domain.ID,
		Price:    domain.Price,
		UserID:   domain.UserID,
		OfficeID: domain.OfficeID,
		CreatedAt: domain.CreatedAt,
		UpdatedAt: domain.UpdatedAt,
		DeletedAt: domain.DeletedAt,
	}
}

func (rec *Transaction) ToDomain() TransactionUseCase.Domain {
	return TransactionUseCase.Domain{
		ID:       rec.ID,
		Price:    rec.Price,
		UserID:   rec.UserID,
		OfficeID: rec.OfficeID,
		CreatedAt: rec.CreatedAt,
		UpdatedAt: rec.UpdatedAt,
		DeletedAt: rec.DeletedAt,
	}
}
