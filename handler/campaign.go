package handler

import (
	"net/http"
	"restful-api/campaign"
	"restful-api/helper"
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