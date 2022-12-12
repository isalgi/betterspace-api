package transactions

import (
	transactions "backend/businesses/transactions"
	// officeRecord "backend/drivers/mysql/offices"
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

	t.conn.Find(&rec)

	TransactionDomain := []transactions.Domain{}

	for _, trans := range rec {
		TransactionDomain = append(TransactionDomain, trans.ToDomain())
	}

	return TransactionDomain
}

func (t *TransactionRepository) GetByUserID(userId string) []transactions.Domain {
	var rec []Transaction

	t.conn.Where("user_id = ?", userId).Find(&rec)

	TransactionDomain := []transactions.Domain{}

	for _, trans := range rec {
		TransactionDomain = append(TransactionDomain, trans.ToDomain())
	}

	return TransactionDomain
}

func (t *TransactionRepository) GetByOfficeID(officeId string) []transactions.Domain {
	var rec []Transaction

	t.conn.Where("office_id = ?", officeId).Find(&rec)

	TransactionDomain := []transactions.Domain{}

	for _, trans := range rec {
		TransactionDomain = append(TransactionDomain, trans.ToDomain())
	}

	return TransactionDomain
}

func (t *TransactionRepository) Create(TransactionDomain *transactions.Domain) transactions.Domain {
	rec := FromDomain(TransactionDomain)

	result := t.conn.Create(&rec)

	result.Last(&rec)

	return rec.ToDomain()
}

func (t *TransactionRepository) GetByID(id string) transactions.Domain {
	var transaction Transaction

	t.conn.First(&transaction, "id = ?", id)

	return transaction.ToDomain()
}

func (t *TransactionRepository) Update(id string, transactionDomain *transactions.Domain) transactions.Domain {
	transaction := t.GetByID(id)

	updatedTransaction := FromDomain(&transaction)
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

	result := t.conn.Where("id = ?", transaction.ID).
		Select("price", "check_in", "check_out", "duration", "payment_method", "status", "drink", "user_id", "office_id").
		Updates(Transaction{
			Price:         updatedTransaction.Price,
			CheckIn:       updatedTransaction.CheckIn,
			CheckOut:      updatedTransaction.CheckOut,
			Duration:      updatedTransaction.Duration,
			PaymentMethod: updatedTransaction.PaymentMethod,
			Status:        updatedTransaction.Status,
			Drink:         updatedTransaction.Drink,
			UserID:        updatedTransaction.UserID,
			OfficeID:      updatedTransaction.OfficeID,
		})

	if result.RowsAffected == 0 {
		updatedTransaction.UserID = 0
		return updatedTransaction.ToDomain()
	}

	return updatedTransaction.ToDomain()
}

func (t *TransactionRepository) Delete(id string) bool {
	var transaction transactions.Domain = t.GetByID(id)

	deletedTransaction := FromDomain(&transaction)

	result := t.conn.Delete(&deletedTransaction)

	return result.RowsAffected != 0
}
