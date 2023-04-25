package services

import (
	"apr-backend/internal/db"
	"apr-backend/internal/model"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type CompanyService interface {
	SaveCompany(com *model.Company) error
	FindCompanies(filter model.CompanyFilter) ([]model.Company, error)
}

const passwordCost = 12

func NewCompanyService(comRepo db.CompanyRepository) CompanyService {
	return companyService{
		comRepo: comRepo,
	}
}

type companyService struct {
	comRepo db.CompanyRepository
}

// FindCompanies implements CompanyService
func (cs companyService) FindCompanies(filter model.CompanyFilter) ([]model.Company, error) {
	return cs.comRepo.FindCompanies(filter)
}

func (cs companyService) SaveCompany(com *model.Company) error {
	pass, err := bcrypt.GenerateFromPassword([]byte(com.Password), passwordCost)
	if err != nil {
		return fmt.Errorf("Error generating password: %w", err)
	}
	com.Password = string(pass)
	return cs.comRepo.SaveCompany(com)
}
