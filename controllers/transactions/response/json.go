package response

import (
	transactions "backend/businesses/transactions"

	"gorm.io/gorm"
)

type Transaction struct {
	ID            uint           `json:"id"`
	CreatedAt     string         `json:"created_at"`
	UpdatedAt     string         `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"deleted_at"`
	Duration      int            `json:"duration"`
	CheckIn       string         `json:"check_in"`
	CheckOut      string         `json:"check_out"`
	Price         uint           `json:"price"`
	Drink         string         `json:"drink"`
	Status        string         `json:"status"`
	PaymentMethod string         `json:"payment_method"`
	UserID        uint           `json:"user_id"`
	UserFullName  string         `json:"user_full_name"`
	UserEmail     string         `json:"user_email"`
	OfficeID      uint           `json:"office_id"`
	OfficeName    string         `json:"office_name"`
}

func FromDomain(domain transactions.Domain) Transaction {
	return Transaction{
		ID:            domain.ID,
		CreatedAt:     domain.CreatedAt.Format("02-01-2006 15:04:05"),
		UpdatedAt:     domain.UpdatedAt.Format("02-01-2006 15:04:05"),
		DeletedAt:     domain.DeletedAt,
		CheckIn:       domain.CheckIn.Format("02-01-2006 15:04:05"),
		CheckOut:      domain.CheckOut.Format("02-01-2006 15:04:05"),
		Duration:      domain.Duration,
		Price:         domain.Price,
		Drink:         domain.Drink,
		Status:        domain.Status,
		PaymentMethod: domain.PaymentMethod,
		UserFullName:  domain.UserFullName,
		UserEmail:     domain.UserEmail,
		OfficeName:    domain.OfficeName,
		UserID:        domain.UserID,
		OfficeID:      domain.OfficeID,
	}
}
