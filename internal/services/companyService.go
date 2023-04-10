package services

import (
	"apr-backend/internal/db"
	"apr-backend/internal/model"
)

type CompanyService interface {
	SaveCompany(com model.Company) error
}

func NewCompanyService(comRepo db.CompanyRepository) CompanyService {
	return companyService{
		comRepo: comRepo,
	}
}

type companyService struct {
	comRepo db.CompanyRepository
}

func (cs companyService) SaveCompany(com model.Company) error {
	return cs.comRepo.SaveCompany(com)
}
