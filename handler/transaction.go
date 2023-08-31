package handler

import (
	"net/http"
	"restful-api/helper"
	"restful-api/transaction"
	"restful-api/user"

	"github.com/gin-gonic/gin"
)

// parameter uri
// ambil paramtete mapping input struct
// call service dengan paramete input struct
// service berbekal campaign id lalu call repo
// repository find data transaction dari campaign id
type transactionHandler struct{
	service transaction.Service
}

func NewTransaction(service transaction.Service) *transactionHandler{
	return &transactionHandler{service}
}

func (h *transactionHandler) GetCampaignTrasaction(c *gin.Context){
	var input transaction.GetCampaignTransactionInput
	err := c.ShouldBindUri(&input)
	if err != nil {
		response := helper.APIResponse("Failed to get campaign's transactions", http.StatusBadRequest, "error",nil)
		c.JSON(http.StatusBadRequest,response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)
	input.User = currentUser


	transactions,err := h.service.GetTransactionByCampaignID(input)
	if err != nil {
		response := helper.APIResponse("Failed to get campaign's transactions", http.StatusBadRequest, "error",nil)
		c.JSON(http.StatusBadRequest,response)
		return
	}
	response:= helper.APIResponse("Campaign's transactions",http.StatusOK,"success",transaction.FormatCampaignTransactions(transactions))
	c.JSON(http.StatusOK,response)
}

// handler
// ambil nilai user dari middlerwar
//service
// repo ambil data transaction dan preload campaign
func(h *transactionHandler) GetUserTransaction(c *gin.Context){
	currentUser := c.MustGet("currentUser").(user.User)
	userID := currentUser.ID

	transactions,err:=h.service.GetTransactionByUserID(userID)
	if err != nil {
		response := helper.APIResponse("Failed to get user's transactions", http.StatusBadRequest, "error",nil)
		c.JSON(http.StatusBadRequest,response)
		return
	}
	response:= helper.APIResponse("User's transactions",http.StatusOK,"success",transaction.FormatUserTransactions(transactions))
	c.JSON(http.StatusOK,response)
}

// ada input dari user 
// handler get input lalu mapping ke struct input
// service menangkap struct  input lalu buat transaksi ,manggil sistem midtrans
// repo create new transaction data

func(h *transactionHandler) CreateTransaction(c *gin.Context){
	var input transaction.CreateTransactionInput

	err := c.ShouldBindJSON(&input)
	
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors":errors}
		response := helper.APIResponse("Failed to create transaction", http.StatusUnprocessableEntity,"error",errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)
	input.User = currentUser

	newTransaction, err := h.service.CreateTransaction(input)
	if err != nil {
		response := helper.APIResponse("Failed to create transaction", http.StatusBadRequest,"error",nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response:= helper.APIResponse("Success create transaction",http.StatusOK,"success",transaction.FormatTransaction(newTransaction))
	c.JSON(http.StatusOK,response)
}