package transaction

import "restful-api/user"

type GetCampaignTransactionInput struct{
	ID int `uri:"id" binding:"required"`
	User user.User
}