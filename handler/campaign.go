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

campaings,err := h.service.GetCampaigns(userID)
if err != nil {

		response := helper.APIResponse("Error to get campaigns", http.StatusBadRequest,"error",nil)
		c.JSON(http.StatusBadRequest, response)
		return
}
response := helper.APIResponse("List of campaigns", http.StatusOK,"success",campaings)
c.JSON(http.StatusOK, response)
}