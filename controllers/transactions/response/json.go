package response

import (
	transactions "backend/businesses/transactions"
)

type Transaction struct {
	ID       uint `json:"id" form:"id"`
	CreatedAt	string		`json:"created_at"`
	UpdatedAt	string		`json:"updated_at"`
	DeletedAt	string	`json:"deleted_at"`
	Price    uint `json:"price" form:"price"`
	UserID   uint `json:"user_id" form:"user_id"`
	OfficeID uint `json:"office_id" form:"office_id"`
}

func FromDomain(domain transactions.Domain) Transaction {
	return Transaction{
		ID:       domain.ID,
		CreatedAt: domain.CreatedAt.Format("02-01-2006 15:04:05"),
		UpdatedAt: domain.UpdatedAt.Format("02-01-2006 15:04:05"),
		DeletedAt: domain.DeletedAt.Time.Format("02-01-2006 15:04:05"),
		Price:    domain.Price,
		UserID:   domain.UserID,
		OfficeID: domain.OfficeID,
	}
}
