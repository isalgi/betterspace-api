package transactions

import (
	transactions "backend/businesses/transactions"
	"fmt"

	"gorm.io/gorm"
)

type TransactionRepository struct {
	conn *gorm.DB
}

func NewMySQLRepository(conn *gorm.DB) transactions.Repository {
	return &TransactionRepository{
		conn: conn,
	}
}

func (t *TransactionRepository) GetAll() []transactions.Domain {
	var rec []Transaction

	t.conn.Preload("User").Preload("Office").Find(&rec)

	TransactionDomain := []transactions.Domain{}

	for _, trans := range rec {
		TransactionDomain = append(TransactionDomain, trans.ToDomain())
	}

	return TransactionDomain
}

func (t *TransactionRepository) GetByUserID(userId string) []transactions.Domain {
	var rec []Transaction

	t.conn.Preload("User").Preload("Office").Where("user_id = ?", userId).Find(&rec)

	TransactionDomain := []transactions.Domain{}

	for _, trans := range rec {
		TransactionDomain = append(TransactionDomain, trans.ToDomain())
	}

	return TransactionDomain
}

func (t *TransactionRepository) GetByOfficeID(officeId string) []transactions.Domain {
	var rec []Transaction

	t.conn.Preload("User").Preload("Office").Where("office_id = ?", officeId).Find(&rec)

	TransactionDomain := []transactions.Domain{}

	for _, trans := range rec {
		TransactionDomain = append(TransactionDomain, trans.ToDomain())
	}

	return TransactionDomain
}

func (t *TransactionRepository) Create(TransactionDomain *transactions.Domain) transactions.Domain {
	rec := FromDomain(TransactionDomain)
	
	result := t.conn.Preload("User").Preload("Office").Create(&rec)

	result.Last(&rec)

	return rec.ToDomain()
}

func (t *TransactionRepository) GetByID(id string) transactions.Domain {
	var transaction Transaction

	t.conn.Preload("User").Preload("Office").First(&transaction, "id = ?", id)

	return transaction.ToDomain()
}

func (t *TransactionRepository) Update(id string, transactionDomain *transactions.Domain) transactions.Domain {
	transaction := t.GetByID(id)
	updatedTransaction := FromDomain(&transaction)

	var getUserID uint
	t.conn.Raw("SELECT ID FROM `users` WHERE `id` = ?", transactionDomain.UserID).Scan(&getUserID)
	fmt.Println("user id", getUserID)
	
	if getUserID == 0 {
		updatedTransaction.ID = 0
		return updatedTransaction.ToDomain()
	}
	
	var getOfficeID uint
	t.conn.Raw("SELECT ID FROM `offices` WHERE `id` = ?", transactionDomain.OfficeID).Scan(&getOfficeID)
	fmt.Println("office id", getOfficeID)
	
	if getOfficeID == 0 {
		updatedTransaction.ID = 0
		return updatedTransaction.ToDomain()
	}

	updatedTransaction.Price = transactionDomain.Price
	updatedTransaction.CheckIn = transactionDomain.CheckIn
	updatedTransaction.CheckOut = transactionDomain.CheckOut
	updatedTransaction.Duration = transactionDomain.Duration
	updatedTransaction.PaymentMethod = transactionDomain.PaymentMethod

	if transactionDomain.Status == "" {
		updatedTransaction.Status = transaction.Status
		fmt.Println(updatedTransaction.Status)
	} else {
		updatedTransaction.Status = transactionDomain.Status
	}

	updatedTransaction.Drink = transactionDomain.Drink
	updatedTransaction.UserID = transactionDomain.UserID
	updatedTransaction.OfficeID = transactionDomain.OfficeID

	t.conn.Preload("User").Preload("Office").Save(&updatedTransaction)

	return updatedTransaction.ToDomain()
}

func (t *TransactionRepository) Delete(id string) bool {
	var transaction transactions.Domain = t.GetByID(id)

	deletedTransaction := FromDomain(&transaction)

	result := t.conn.Delete(&deletedTransaction)

	return result.RowsAffected != 0
}
