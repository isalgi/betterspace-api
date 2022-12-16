package drivers

import (
	facilityDomain "backend/businesses/facilities"
	officeFacilityDomain "backend/businesses/office_facilities"
	officeImageDomain "backend/businesses/office_images"
	officeDomain "backend/businesses/offices"
	reviewDomain "backend/businesses/review"
	transactionDomain "backend/businesses/transactions"
	userDomain "backend/businesses/users"

	facilityDB "backend/drivers/mysql/facilities"
	officeFacilityDB "backend/drivers/mysql/office_facilities"
	officeImageDB "backend/drivers/mysql/office_images"
	officeDB "backend/drivers/mysql/offices"
	reviewDB "backend/drivers/mysql/review"
	transactionDB "backend/drivers/mysql/transactions"
	userDB "backend/drivers/mysql/users"

	"gorm.io/gorm"
)

func NewUserRepository(conn *gorm.DB) userDomain.Repository {
	return userDB.NewMySQLRepository(conn)
}

func NewOfficeRepository(conn *gorm.DB) officeDomain.Repository {
	return officeDB.NewMySQLRepository(conn)
}

func NewOfficeImageRepository(conn *gorm.DB) officeImageDomain.Repository {
	return officeImageDB.NewMySQLRepository(conn)
}

func NewFacilityRepository(conn *gorm.DB) facilityDomain.Repository {
	return facilityDB.NewMySQLRepository(conn)
}

func NewOfficeFacilityRepository(conn *gorm.DB) officeFacilityDomain.Repository {
	return officeFacilityDB.NewMySQLRepository(conn)
}

func NewTransactionRepository(conn *gorm.DB) transactionDomain.Repository {
	return transactionDB.NewMySQLRepository(conn)
}

func NewReviewRepository(conn *gorm.DB) reviewDomain.Repository {
	return reviewDB.NewMySQLRepository(conn)
}
