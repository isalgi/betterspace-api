package request

import (
	transactions "backend/businesses/transactions"
	"backend/utils"
	"time"

	"github.com/go-playground/validator/v10"
)

type Transaction struct {
	Price    uint `json:"price" form:"price" validate:"required"`
	UserID   uint `json:"user_id" form:"user_id" validate:"required"`
	CheckIn  time.Time
	Drink string `json:"drink" form:"drink" validate:"required"`
	Duration int `json:"duration" form:"duration" validate:"required"`
	OfficeID uint `json:"office_id" form:"office_id" validate:"required"`
}

type CheckInDTO struct {
	CheckInHour string `json:"check_in_hour" form:"check_in_hour"`
	CheckInDate string `json:"check_in_date" form:"check_in_date"`
}

func (req *Transaction) ToDomain() *transactions.Domain {
	return &transactions.Domain{
		Price:    req.Price,
		CheckIn: req.CheckIn,
		Duration: req.Duration,
		Drink: req.Drink,
		UserID:   req.UserID,
		OfficeID: req.OfficeID,
	}
}

func (req *Transaction) Validate() error {
	validate := validator.New()

	err := validate.Struct(req)

	return err
}

func (req *CheckInDTO) Validate() error {
	err := utils.IsValidTime(req.CheckInHour)

	if err != nil {
		return err
	}

	err = utils.DateValidation(req.CheckInDate)

	return err
}