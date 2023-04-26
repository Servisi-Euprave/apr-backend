package services

import (
	"apr-backend/internal/db"
	"apr-backend/internal/model"
)

type NstjService interface {
	FindAll() ([]model.Nstj, error)
}

type nstjService struct {
	nstjRepo db.NstjRepository
}

// FindAll implements NstjService
func (ns nstjService) FindAll() ([]model.Nstj, error) {
	return ns.nstjRepo.FindAll()
}

func NewNstjService(nstjRepo db.NstjRepository) NstjService {
	return nstjService{
		nstjRepo: nstjRepo,
	}
}
