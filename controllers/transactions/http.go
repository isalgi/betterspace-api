package transactions

import (
	transactions "backend/businesses/transactions"
	ctrl "backend/controllers"
	"backend/controllers/transactions/request"
	"backend/controllers/transactions/response"
	"backend/utils"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type TransactionController struct {
	TransactionUsecase transactions.Usecase
}

func NewTransactionController(tc transactions.Usecase) *TransactionController {
	return &TransactionController{
		TransactionUsecase: tc,
	}
}

func (t *TransactionController) GetAll(c echo.Context) error {
	TransactionsData := t.TransactionUsecase.GetAll()

	Transactions := []response.Transaction{}

	for _, trans := range TransactionsData {
		Transactions = append(Transactions, response.FromDomain(trans))
	}

	return ctrl.NewResponse(c, http.StatusOK, "success", "all transactions", Transactions)
}

func (t *TransactionController) Create(c echo.Context) error {
	input := request.Transaction{}

	if err := c.Bind(&input); err != nil {
		return ctrl.NewInfoResponse(c, http.StatusBadRequest, "failed", "validation failed")
	}

	checkInDTO := request.CheckInDTO{}

	if err := c.Bind(&checkInDTO); err != nil {
		return ctrl.NewInfoResponse(c, http.StatusBadRequest, "failed", "bind time failed")
	}

	// input hour validation
	if err := checkInDTO.Validate(); err != nil {
		return ctrl.NewInfoResponse(c, http.StatusBadRequest, "failed", fmt.Sprintf("%s", err))
	}

	checkInHour := utils.ConvertStringToShiftTime(checkInDTO.CheckInDate, checkInDTO.CheckInHour)
	
	input.CheckIn = checkInHour

	err := input.Validate()

	if err != nil {
		return ctrl.NewInfoResponse(c, http.StatusBadRequest, "failed", "validation failed")
	}

	trans := t.TransactionUsecase.Create(input.ToDomain())

	if trans.ID == 0 {
		return ctrl.NewInfoResponse(c, http.StatusBadRequest, "failed", "create transaction failed")
	}

	return ctrl.NewResponse(c, http.StatusCreated, "success", "transaction created", response.FromDomain(trans))
}

func (t *TransactionController) GetByID(c echo.Context) error {
	var id string = c.Param("id")

	transaction := t.TransactionUsecase.GetByID(id)

	if transaction.ID == 0 {
		return ctrl.NewResponse(c, http.StatusNotFound, "failed", "transaction not found", "")
	}

	return ctrl.NewResponse(c, http.StatusOK, "success", "transaction found", response.FromDomain(transaction))
}

func (t *TransactionController) Update(c echo.Context) error {
	input := request.Transaction{}

	if err := c.Bind(&input); err != nil {
		return ctrl.NewResponse(c, http.StatusBadRequest, "failed", "validation failed", "")
	}

	var transactionId string = c.Param("id")

	err := input.Validate()

	if err != nil {
		return ctrl.NewResponse(c, http.StatusBadRequest, "failed", "validation failed", "")
	}

	transaction := t.TransactionUsecase.Update(transactionId, input.ToDomain())

	if transaction.ID == 0 {
		return ctrl.NewResponse(c, http.StatusNotFound, "failed", "transaction not found", "")
	}

	return ctrl.NewResponse(c, http.StatusOK, "success", "transaction updated", response.FromDomain(transaction))
}

func (t *TransactionController) Delete(c echo.Context) error {
	var transactionId string = c.Param("id")

	isSuccess := t.TransactionUsecase.Delete(transactionId)

	if !isSuccess {
		return ctrl.NewResponse(c, http.StatusNotFound, "failed", "transaction not found", "")
	}

	return ctrl.NewResponse(c, http.StatusOK, "success", "transaction deleted", "")
}
