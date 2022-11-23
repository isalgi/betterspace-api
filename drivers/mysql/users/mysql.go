package users

import (
	"backend/businesses/users"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type userRepository struct {
	conn *gorm.DB
}

func NewMySQLRepository(conn *gorm.DB) users.Repository {
	return &userRepository{
		conn: conn,
	}
}

func (ur *userRepository) Register(userDomain *users.Domain) users.Domain {
	password, _ := bcrypt.GenerateFromPassword([]byte(userDomain.Password), bcrypt.DefaultCost)

	rec := FromDomain(userDomain)

	uuid := uuid.New().String()
	uuidWithoutHyphens := strings.Replace(uuid, "-", "", -1)

	rec.ID = uuidWithoutHyphens
	rec.Password = string(password)
	rec.Image = ""
	rec.Roles = "user"

	var user User
	ur.conn.First(&user, "email = ?", userDomain.Email)

	// handle email if email already used for an account
	if user.ID != "" {
		fmt.Println("user exist. proceed to login or use another email.")
		return users.Domain{}
	}

	result := ur.conn.Create(&rec)
	result.Last(&rec)

	return rec.ToDomain()
}

func (ur *userRepository) GetByEmail(userDomain *users.LoginDomain) users.Domain {
	var user User
	ur.conn.First(&user, "email = ?", userDomain.Email)

	if user.ID == "" {
		fmt.Println("user not found")
		return users.Domain{}
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userDomain.Password))

	if err != nil {
		fmt.Println("password failed")
		return users.Domain{}
	}

	return user.ToDomain()
}