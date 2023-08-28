package transaction

import (
	"errors"
	"restful-api/campaign"
)

type service struct{
	repository Repository
	campaignRepository campaign.Repository
}

type Service interface{
	GetTransactionByCampaignID(input GetCampaignTransactionInput)([]Transaction,error)
}
func NewService(repository Repository,campaignRepsitory campaign.Repository) *service{
	return &service{repository,campaignRepsitory}
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