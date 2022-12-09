package request

import (
	transactions "backend/businesses/transactions"

	"github.com/go-playground/validator/v10"
)

type Transaction struct {
	Price    uint `json:"price" form:"price" validate:"required"`
	UserID   uint `json:"user_id" form:"user_id" validate:"required"`
	OfficeID uint `json:"office_id" form:"office_id" validate:"required"`
}

func (req *Transaction) ToDomain() *transactions.Domain {
	return &transactions.Domain{
		Price:    req.Price,
		UserID:   req.UserID,
		OfficeID: req.OfficeID,
	}
}

func (req *Transaction) Validate() error {
	validate := validator.New()

	err := validate.Struct(req)

	return err
}
