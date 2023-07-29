package handler

import (
	"net/http"
	"restful-api/helper"
	"restful-api/user"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService user.Service
}
func NewUserHandler(userService user.Service) *userHandler{
	return &userHandler{userService}
}

func (h *userHandler)RegisterUser(c *gin.Context){
//catch json input from user
//map the input to RegisterUserInput
//Passing input to service

	var input user.RegisterUserInput
	err := c.ShouldBindJSON(&input)
	if err != nil{
		errors:=helper.FormatValidationError(err)
		errorMessage := gin.H{"errors":errors}

		response := helper.APIResponse("Register account failed", http.StatusUnprocessableEntity,"error",errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	newUser,err := h.userService.RegisterUser(input)
	if err!=nil{
		response := helper.APIResponse("Register account failed", http.StatusBadRequest,"error",nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	formatter := user.FormatJSONUser(newUser,"tokenjwtcontoh")
	response := helper.APIResponse("Account has been created", http.StatusOK,"success",formatter)

	c.JSON(http.StatusOK,response)
}

func (h *userHandler) LoginUser(c *gin.Context){
	//catch user input by handler
	var input user.LoginInput
	err := c.ShouldBindJSON(&input)
	if err != nil{
		errors:=helper.FormatValidationError(err)
		errorMessage := gin.H{"errors":errors}

		response := helper.APIResponse("Login failed", http.StatusUnprocessableEntity,"error",errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	LoggedinUser, err := h.userService.LoginUser(input)
	if err != nil{
		errorMessage := gin.H{"errors":err.Error()}

		response := helper.APIResponse("Login failed", http.StatusUnprocessableEntity,"error",errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	formatter := user.FormatJSONUser(LoggedinUser,"token")
	response := helper.APIResponse("Successfully Logged in", http.StatusOK,"success",formatter)

	c.JSON(http.StatusOK,response)
	//mapping user input into struct LoginInputUser
	//passing input struct to service
	//service will match the user input in db by his email
	//matching the password after match the email is true
}