package handler

import (
	"fmt"
	"net/http"
	"restful-api/auth"
	"restful-api/helper"
	"restful-api/user"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService user.Service
	authService auth.Service
}
func NewUserHandler(userService user.Service,authService auth.Service) *userHandler{
	return &userHandler{userService,authService}
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
	token,err:=h.authService.GenerateToken(newUser.ID)
	if err!=nil{
		response := helper.APIResponse("Register account failed", http.StatusBadRequest,"error",nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	formatter := user.FormatJSONUser(newUser,token)
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
	token,err:=h.authService.GenerateToken(LoggedinUser.ID )
	if err!=nil{
		response := helper.APIResponse("Login failed", http.StatusBadRequest,"error",nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	formatter := user.FormatJSONUser(LoggedinUser,token)
	response := helper.APIResponse("Successfully Logged in", http.StatusOK,"success",formatter)

	c.JSON(http.StatusOK,response)
	//mapping user input into struct LoginInputUser
	//passing input struct to service
	//service will match the user input in db by his email
	//matching the password after match the email is true
}

func (h *userHandler) CheckEmailAvailability (c *gin.Context){
	// input email user
	var input user.CheckEmailInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors:=helper.FormatValidationError(err)
		errorMessage := gin.H{"errors":errors}

		response := helper.APIResponse("Email Checking Failed", http.StatusUnprocessableEntity,"error",errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	isEmailAvailable, err := h.userService.IsEmailAvailable(input)
	if err != nil {
		errorMessage := gin.H{"errors":"Server Error"}

		response := helper.APIResponse("Email Checking Failed", http.StatusUnprocessableEntity,"error",errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	data := gin.H{
		"is_available" : isEmailAvailable,
	}
	metaMessage := "Email has been registered"
	if isEmailAvailable{
		metaMessage = "Email is Available"
	}
	response := helper.APIResponse(metaMessage, http.StatusOK,"success",data)
	c.JSON(http.StatusOK, response)
		
	// mapping input email ke struct Login input
	// struct di passing ke service
	// servicer manggil repository untuk cek email sudah ada /belum
	//repositor db
}

func (h *userHandler) UploadAvatar (c *gin.Context){
file,err:= c.FormFile("avatar")
if err != nil{
	data := gin.H{"is_uploaded":false}
	response := helper.APIResponse("Failed upload avatar image",http.StatusBadRequest,"error",data)
	c.JSON(http.StatusBadRequest,response)
	return
}
//dapet dari JWT
userID := 1

path := fmt.Sprintf("images/%d-%s" ,userID, file.Filename) 
err = c.SaveUploadedFile(file,path)
if err != nil{
	data := gin.H{"is_uploaded":false}
	response := helper.APIResponse("Failed upload avatar image",http.StatusBadRequest,"error",data)
	c.JSON(http.StatusBadRequest,response)
	return
}

_,err =h.userService.SaveAvatar(userID,path)
if err != nil{
	data := gin.H{"is_uploaded":false}
	response := helper.APIResponse("Failed upload avatar image",http.StatusBadRequest,"error",data)
	c.JSON(http.StatusBadRequest,response)
	return
}
data := gin.H{"is_uploaded":true}
	response := helper.APIResponse("Avatar successfully update",http.StatusOK,"success",data)
	c.JSON(http.StatusOK,response)

//catch input user form body
	//save gambar di folder "images/"
	//service memanggil repo
	//JWT untuk memperoleh id (hrdcode = ID 1)
	//repo get user id = 1
	//repo update data user simpan lokasi file
}