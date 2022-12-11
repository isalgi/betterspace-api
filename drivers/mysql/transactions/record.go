package transactions

import (
	TransactionUseCase "backend/businesses/transactions"
	"backend/drivers/mysql/offices"
	"backend/drivers/mysql/users"
	"time"

	"gorm.io/gorm"
)

type Transaction struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
	Price     uint           `json:"price"`
	CheckIn   time.Time      `json:"check_in" gorm:"type:timestamp;not null;default:now()"`
	Duration  int            `json:"duration" form:"duration"`
	CheckOut  time.Time      `json:"check_out" gorm:"type:timestamp;not null;default:now()"`
	Drink     string         `json:"drink" form:"drink"`
	UserID    uint           `json:"user_id"`
	OfficeID  uint           `json:"office_id"`
	User      users.User     `json:"user" gorm:"foreignKey:UserID;references:ID"`
	Office    offices.Office `json:"office" gorm:"foreignKey:OfficeID;references:ID"`
}

func FromDomain(domain *TransactionUseCase.Domain) *Transaction {
	return &Transaction{
		ID:        domain.ID,
		Price:     domain.Price,
		CheckIn:   domain.CheckIn,
		Duration:  domain.Duration,
		CheckOut:  domain.CheckOut,
		Drink:     domain.Drink,
		UserID:    domain.UserID,
		OfficeID:  domain.OfficeID,
		CreatedAt: domain.CreatedAt,
		UpdatedAt: domain.UpdatedAt,
		DeletedAt: domain.DeletedAt,
	}
}

func (rec *Transaction) ToDomain() TransactionUseCase.Domain {
	return TransactionUseCase.Domain{
		ID:        rec.ID,
		Price:     rec.Price,
		CheckIn:   rec.CheckIn,
		Duration:  rec.Duration,
		CheckOut:  rec.CheckOut,
		Drink:     rec.Drink,
		UserID:    rec.UserID,
		OfficeID:  rec.OfficeID,
		CreatedAt: rec.CreatedAt,
		UpdatedAt: rec.UpdatedAt,
		DeletedAt: rec.DeletedAt,
	}
}
