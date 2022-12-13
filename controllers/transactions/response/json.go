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
	User          struct {
		UserID   uint   `json:"user_id"`
		FullName string `json:"full_name"`
		Email    string `json:"email"`
	}
	Office struct {
		OfficeID   uint   `json:"office_id"`
		OfficeName string `json:"office_name"`
	}
}

type user struct {
	UserID   uint   `json:"user_id"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
}

type office struct {
	OfficeID   uint   `json:"office_id"`
	OfficeName string `json:"office_name"`
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
		User: user{
			UserID:   domain.UserID,
			FullName: domain.UserFullName,
			Email:    domain.UserEmail,
		},
		Office: office{
			OfficeID:   domain.OfficeID,
			OfficeName: domain.OfficeName,
		},
	}
}
