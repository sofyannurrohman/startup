package handler

import (
	"fmt"
	"net/http"
	"restful-api/campaign"
	"restful-api/helper"
	"restful-api/user"
	"strconv"

	"github.com/gin-gonic/gin"
)


type campaignHandler struct{
	service campaign.Service
}
func NewCampaignHandler(service campaign.Service) *campaignHandler{
	return &campaignHandler{service}
}


func(h *campaignHandler) GetCampaigns(c *gin.Context){

userID,_ := strconv.Atoi(c.Query("user_id"))

campaigns,err := h.service.GetCampaigns(userID)
if err != nil {

		response := helper.APIResponse("Error to get campaigns", http.StatusBadRequest,"error",nil)
		c.JSON(http.StatusBadRequest, response)
		return
}
response := helper.APIResponse("List of campaigns", http.StatusOK,"success",campaign.FormatCampaigns(campaigns))
c.JSON(http.StatusOK, response)
}

func(h *campaignHandler) GetCampaign(c *gin.Context){
	//api/v1.campaigns/2
	//handler mapping id di yg ada di url ke struct input => service, dan call formatter
	// service dari struct input menangkap id di url lalu call repo
	//repository get campaign by id
	var input campaign.GetCampaignDetailInput
	err:=c.ShouldBindUri(&input)
	if err != nil{
		response := helper.APIResponse("Failed to get detail of campaign", http.StatusBadRequest, "error",nil)
		c.JSON(http.StatusBadRequest,response)
		return
	}
	campaignDetail,err := h.service.GetCampaignByID(input)
	if err != nil{
		response := helper.APIResponse("Failed to get detail of campaign", http.StatusBadRequest, "error",nil)
		c.JSON(http.StatusBadRequest,response)
		return
	}
	response:= helper.APIResponse("Campaign detail",http.StatusOK,"success",campaign.FormatCampaignDetail(campaignDetail))
	c.JSON(http.StatusOK,response)
}

//tangkap parameter dari user lalu mapping ke struct input
//ambil current user dari jwt
func(h *campaignHandler) CreateCampaign(c *gin.Context){
	var input campaign.CreateCampaignInput
	err:=c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors":errors}
		response := helper.APIResponse("Failed to create campaign", http.StatusUnprocessableEntity,"error",errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	currentUser := c.MustGet("currentUser").(user.User)
	input.User = currentUser
	newCampaign,err := h.service.CreateCampaign(input)
	if err != nil {

		response := helper.APIResponse("Failed to create campaign", http.StatusBadRequest,"error",nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response:= helper.APIResponse("Success to create campaign",http.StatusOK,"success",campaign.FormatCampaign(newCampaign))
	c.JSON(http.StatusOK,response)
}
// handler menangkap input dari user
func(h *campaignHandler) UpdateCampaign(c *gin.Context){
	var inputID campaign.GetCampaignDetailInput
	err :=c.ShouldBindUri(&inputID)
	if err != nil{
		response := helper.APIResponse("Failed to update campaign", http.StatusBadRequest, "error",nil)
		c.JSON(http.StatusBadRequest,response)
		return
	}

	var inputData campaign.CreateCampaignInput
	currentUser := c.MustGet("currentUser").(user.User)
	inputData.User = currentUser

	err =c.ShouldBindJSON(&inputData)
	if err != nil {

		response := helper.APIResponse("Failed to update campaign", http.StatusUnprocessableEntity,"error",nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	updatedCampaign,err:= h.service.UpdateCampaign(inputID,inputData)
	if err != nil {
		response := helper.APIResponse("Failed to update campaign", http.StatusUnprocessableEntity,"error",nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response:= helper.APIResponse("Success to update campaign",http.StatusOK,"success",campaign.FormatCampaign(updatedCampaign))
	c.JSON(http.StatusOK,response)
}

func(h *campaignHandler) UploadImage(c *gin.Context){
	var input campaign.CreateCampaignImageInput
	err := c.ShouldBind(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors":errors}
		response := helper.APIResponse("Failed to upload campaign image", http.StatusBadRequest,"error",errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)
	input.User = currentUser
	userID := currentUser.ID
	file,err:= c.FormFile("file")
if err != nil{
	data := gin.H{"is_uploaded":false}
	response := helper.APIResponse("Failed upload campaign image",http.StatusBadRequest,"error",data)
	c.JSON(http.StatusBadRequest,response)
	return
}
//dapet dari JWT

path := fmt.Sprintf("images/%d-%s" ,userID, file.Filename) 
err = c.SaveUploadedFile(file,path)
if err != nil{
	data := gin.H{"is_uploaded":false}
	response := helper.APIResponse("Failed upload campaign image",http.StatusBadRequest,"error",data)
	c.JSON(http.StatusBadRequest,response)
	return
}
_,err = h.service.SaveCampaignImage(input,path)
if err != nil{
	data := gin.H{"is_uploaded":false}
	response := helper.APIResponse("Failed upload campaign image",http.StatusBadRequest,"error",data)
	c.JSON(http.StatusBadRequest,response)
	return
}
data := gin.H{"is_uploaded":true}
	response := helper.APIResponse("Campaign successfully update",http.StatusOK,"success",data)
	c.JSON(http.StatusOK,response)

}
// handler 
// tangkap input dan ubah ke struct input
// save image ke folder
// service call repository
// repository :
// 1. save data image
// 2. ubah is_primary true ke false ketika upload is_primary double