package transactions

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
