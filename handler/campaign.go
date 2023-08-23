package handler

import (
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

		response := helper.APIResponse("Failed to create campaign", http.StatusUnprocessableEntity,"error",nil)
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