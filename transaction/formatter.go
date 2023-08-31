package transaction

import (
	"time"
)

type CampaignTransactionFormatter struct{
	ID int `json:"id"`
	Name string `json:"name"`
	Amount int `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func FormatCampaignTransaction (transaction Transaction) CampaignTransactionFormatter{
	formatter:= CampaignTransactionFormatter{}
	formatter.ID = transaction.ID
	formatter.Name = transaction.User.Name
	formatter.Amount =transaction.Amount
	formatter.CreatedAt = transaction.CreatedAt
	formatter.UpdatedAt = transaction.UpdatedAt
	return formatter

}
func FormatCampaignTransactions (transaction []Transaction) []CampaignTransactionFormatter{
	if len(transaction) == 0 {
		return []CampaignTransactionFormatter{}
	}
	var transactionFormatter []CampaignTransactionFormatter
	for _,transaction := range transaction{
		formatter := FormatCampaignTransaction(transaction)
		transactionFormatter = append(transactionFormatter, formatter)
	}
	return transactionFormatter
}

type UserTransactionFormatter struct{
	ID int `json:"id"`
	Amount int `json:"ammount"`
	Status string `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	Campaign CampaignFormatter `json:"campaign"`
}
type CampaignFormatter struct {
	Name string `json:"name"`
	ImageUrl string `json:"image_url"`
}

func FormatUserTransaction (transaction Transaction)UserTransactionFormatter{
	formatter := UserTransactionFormatter{}
	formatter.ID = transaction.ID
	formatter.Amount = transaction.Amount
	formatter.Status = transaction.Status
	formatter.CreatedAt = transaction.CreatedAt

	campaignFormater := CampaignFormatter{}
	campaignFormater.Name = formatter.Campaign.Name
	campaignFormater.ImageUrl = ""
	if len (transaction.Campaign.CampaignImages)>0 {
		campaignFormater.ImageUrl = transaction.Campaign.CampaignImages[0].FileName
	}
	
	formatter.Campaign = campaignFormater
	return formatter
}

func FormatUserTransactions (transactions []Transaction)[]UserTransactionFormatter{
	if len(transactions) == 0 {
		return []UserTransactionFormatter{}
	}
	var transactionsFormatter []UserTransactionFormatter
	for _,transactions := range transactions{
		formatter := FormatUserTransaction(transactions)
		transactionsFormatter = append(transactionsFormatter, formatter)
	}
	return transactionsFormatter
}

type TransactionFormatter struct{
	ID int `json:"id"`
	CampaignID int  `json:"campaign_id"`
	UserID int `json:"user_id"`
	Amount int `json:"amount"`
	Status string `json:"status"`
	Code string `json:"code"`
	PaymentURL string `json:"payment_url"`
}

func FormatTransaction (transaction Transaction) TransactionFormatter{
	formatter := TransactionFormatter{}
	formatter.ID = transaction.ID
	formatter.CampaignID = transaction.CampaignID
	formatter.UserID = transaction.UserID
	formatter.Amount = transaction.Amount
	formatter.Status = transaction.Status
	formatter.Code = transaction.Code
	formatter.PaymentURL = transaction.PaymentURL

	return formatter
}