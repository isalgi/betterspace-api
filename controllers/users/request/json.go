package request

import (
	"backend/businesses/users"

	"github.com/go-playground/validator/v10"
)

type User struct {
	FullName				string `json:"full_name" validate:"required"`
	Gender					string `json:"gender" validate:"required"`
	Email					string `json:"email" validate:"required,email"`
	Password				string `json:"password" validate:"required"`
	ConfirmationPassword	string `json:"confirmation_password" validate:"required"`
	Image 					string `json:"image" form:"image"`
	Roles 					string `json:"roles"`
}

type UserLogin struct {
	Email					string `json:"email" validate:"required,email"`
	Password				string `json:"password" validate:"required"`
}

func (req *User) ToDomainRegister() *users.Domain {
	return &users.Domain{
		FullName: req.FullName,
		Gender: req.Gender,
		Email:    req.Email,
		Password: req.Password,
	}
}

func (req *UserLogin) ToDomainLogin() *users.LoginDomain {
	return &users.LoginDomain{
		Email: req.Email,
		Password: req.Password,
	}
}

func (req *User) Validate() error {
	validate := validator.New()
	
	err := validate.Struct(req)

	return err
}

func (req *UserLogin) Validate() error {
	validate := validator.New()
	
	err := validate.Struct(req)

	return err
}