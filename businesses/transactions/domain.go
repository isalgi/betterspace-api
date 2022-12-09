package transactions

type Domain struct {
	ID       uint
	Price    uint
	UserID   uint
	OfficeID uint
}

type Usecase interface {
	GetAll() []Domain
	Create(transactionDomain *Domain) Domain
	GetByID(id string) Domain
	Update(id string, transactionDomain *Domain) Domain
	Delete(id string) bool
}

type Repository interface {
	GetAll() []Domain
	Create(transactionDomain *Domain) Domain
	GetByID(id string) Domain
	Update(id string, transactionDomain *Domain) Domain
	Delete(id string) bool
}
