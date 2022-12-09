package transactions

import (
	transactions "backend/businesses/transactions"

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
	updatedTransaction.UserID = transactionDomain.UserID
	updatedTransaction.OfficeID = transactionDomain.OfficeID

	t.conn.Where("id = ?", transaction.ID).Select("price", "user_id", "office_id").Updates(Transaction{Price: transactionDomain.Price, UserID: transactionDomain.UserID, OfficeID: transactionDomain.OfficeID})

	return updatedTransaction.ToDomain()
}

func (t *TransactionRepository) Delete(id string) bool {
	var transaction transactions.Domain = t.GetByID(id)

	deletedTransaction := FromDomain(&transaction)

	result := t.conn.Delete(&deletedTransaction)

	return result.RowsAffected != 0
}
