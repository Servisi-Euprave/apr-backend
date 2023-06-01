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
	FindOne(pib int) (model.Company, error)
	LiquidateById(pib int) error
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

// LiquidateById implements CompanyService
func (cs companyService) LiquidateById(pib int) error {
	return cs.comRepo.LiquidateById(pib)
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

func (cs companyService) FindOne(pib int) (model.Company, error) {
	return cs.comRepo.FindOne(pib)
}
