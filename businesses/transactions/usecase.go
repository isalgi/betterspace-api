package transactions

import (
	"time"
)

type transactionUsecase struct {
	transactionRepository Repository
}

func NewTransactionUsecase(tr Repository) Usecase {
	return &transactionUsecase{
		transactionRepository: tr,
	}
}

func (tu *transactionUsecase) GetAll() []Domain {
	return tu.transactionRepository.GetAll()
}

func (tu *transactionUsecase) Create(transactionDomain *Domain) Domain {
	hour, min, sec := transactionDomain.CheckIn.Clock()
	year, month, day := transactionDomain.CheckIn.Date()

	loc, err := time.LoadLocation("Asia/Jakarta")

	if err != nil {
		transactionDomain.ID = 0
		return *transactionDomain
	}

	if transactionDomain.Duration > 13 {
		transactionDomain.ID = 0
		return *transactionDomain
	}

	duration := transactionDomain.Duration
	hour += hour + duration

	if hour > 24 {
		hour %= 24
	}

	transactionDomain.CheckOut = time.Date(year, month, day, hour, min, sec, 0, loc)

	if transactionDomain.CheckIn.After(transactionDomain.CheckOut) {
		transactionDomain.CheckOut = transactionDomain.CheckOut.AddDate(0, 0, 1)
	}

	return tu.transactionRepository.Create(transactionDomain)
}

func (tu *transactionUsecase) GetByID(id string) Domain {
	return tu.transactionRepository.GetByID(id)
}

func (tu *transactionUsecase) Update(id string, transactionDomain *Domain) Domain {
	return tu.transactionRepository.Update(id, transactionDomain)
}

func (tu *transactionUsecase) Delete(id string) bool {
	return tu.transactionRepository.Delete(id)
}
