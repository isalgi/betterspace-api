package response

import transactions "backend/businesses/transactions"

type Transaction struct {
	ID       uint `json:"id" form:"id"`
	Price    uint `json:"price" form:"price"`
	UserID   uint `json:"user_id" form:"user_id"`
	OfficeID uint `json:"office_id" form:"office_id"`
}

func FromDomain(domain transactions.Domain) Transaction {
	return Transaction{
		ID:       domain.ID,
		Price:    domain.Price,
		UserID:   domain.UserID,
		OfficeID: domain.OfficeID,
	}
}
