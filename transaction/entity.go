package transaction

import (
	"restful-api/campaign"
	"restful-api/user"
	"time"
)

type Transaction struct{
	ID int
	CampaignID int
	UserID int
	Amount int
	Status string
	Code string
	User user.User
	PaymentURL string
	Campaign campaign.Campaign
	CreatedAt time.Time
	UpdatedAt time.Time
}