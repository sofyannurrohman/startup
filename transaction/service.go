package transaction

import (
	"errors"
	"restful-api/campaign"
	"restful-api/payment"
	"strconv"
)

type service struct{
	repository Repository
	campaignRepository campaign.Repository
	paymentService payment.Service
}

type Service interface{
	GetTransactionByCampaignID(input GetCampaignTransactionInput)([]Transaction,error)
	GetTransactionByUserID(userID int )([]Transaction,error)
	CreateTransaction(input CreateTransactionInput)(Transaction,error)
	ProccessPayment(input TransactionNotificationInput)(error)
}
func NewService(repository Repository,campaignRepsitory campaign.Repository,paymentService payment.Service) *service{
	return &service{repository,campaignRepsitory,paymentService}
}
func (s *service) GetTransactionByCampaignID(input GetCampaignTransactionInput)([]Transaction,error){
	// get campaign 
	campaign,err := s.campaignRepository.FindByID(input.ID)
	if err != nil{
		return []Transaction{},err
	}
	if campaign.UserID != input.User.ID{
		return []Transaction{},errors.New("not an Owner of the campaign")
	}
	// check campaign id
	transactions,err := s.repository.GetByCampaignID(input.ID)
	if err != nil{
		return transactions,err
	}
	return transactions,nil
}

func (s *service)GetTransactionByUserID(userID int )([]Transaction,error){
	transactions,err := s.repository.GetByUserID(userID)
	if err != nil{
		return transactions,err
	}
	return transactions,nil
}

func (s *service)CreateTransaction(input CreateTransactionInput)(Transaction,error){
	transaction := Transaction{}
	transaction.CampaignID = input.CampaignID
	transaction.Amount = input.Amount
	transaction.UserID = input.User.ID
	transaction.Status = "pending"

	newTransaction,err:=s.repository.Save(transaction)
	if err != nil{
		return newTransaction,err
	}
	paymentTransaction := payment.Transaction{
		ID: newTransaction.ID,
		Amount: newTransaction.Amount,
	}
	paymentURL,err := s.paymentService.GetPaymentURL(paymentTransaction,input.User)
	if err != nil{
		return newTransaction,err
	}
	newTransaction.PaymentURL = paymentURL
 	newTransaction,err = s.repository.Update(newTransaction)
	 if err != nil{
		return newTransaction,err
	}
	return newTransaction,nil
}

func (s *service) ProccessPayment(input TransactionNotificationInput)(error){
	transaction_id,_ := strconv.Atoi(input.OrderID)
	//find transaction
	transaction,err := s.repository.GetByID(transaction_id)
	if err != nil {
		return err
	}

	if (input.PaymentType == "credit_card" && input.TransactionStatus == "capture" && input.FraudStatus=="accept"){
		transaction.Status = "paid"
	}else if input.TransactionStatus == "settlement"{
		transaction.Status = "paid"
	}else if input.TransactionStatus == "deny" || input.TransactionStatus == "expire" || input.TransactionStatus == "cancel"{
		transaction.Status = "cancelled"
	}
	
	updatedTransaction,err:= s.repository.Update(transaction)
	if err != nil {
		return err
	}

	campaign, err:=s.campaignRepository.FindByID(transaction.CampaignID)
	if err != nil {
		return err
	}
	if updatedTransaction.Status=="paid"{
		campaign.BackerCount = campaign.BackerCount + 1
		campaign.CurrentAmount = campaign.CurrentAmount + updatedTransaction.Amount
		_,err:=s.campaignRepository.Update(campaign)
		if err != nil {
			return err
		}	
	}

	return nil
}