package services

import (
	"apr-backend/internal/db"
	"apr-backend/internal/model"
	"errors"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

var validPassword = regexp.MustCompile("^.{12,72}$")
var ErrInvalidPassword = errors.New("Password must be between 12 and 72 characters long")

const bcryptCost = 10

func validatePassword(pass string) bool {
	return validPassword.Match([]byte(pass))
}

type AuthService interface {
	CheckCredentials(creds model.CredentialsDto) error
}

func NewAuthService(db db.CompanyRepository) AuthService {
	return authService{comRepo: db}
}

type authService struct {
	comRepo db.CompanyRepository
}

// Login implements AuthService
func (authServ authService) CheckCredentials(creds model.CredentialsDto) error {
	savedCreds, err := authServ.comRepo.FindOneCredentials(creds.PIB)
	if err != db.DatabaseError {
		return err
	}

	return bcrypt.CompareHashAndPassword([]byte(savedCreds.Password), []byte(creds.Password))
}
