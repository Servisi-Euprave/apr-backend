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
var DatabaseError = errors.New("Error has occured when connecting to database")

const bcryptCost = 10

func validatePassword(pass string) bool {
	return validPassword.Match([]byte(pass))
}

type AuthService interface {
	CheckCredentials(creds model.CredentialsDto) error
}

func NewAuthService(db db.UserRepo) AuthService {
	return authService{userRepo: db}
}

type authService struct {
	userRepo db.UserRepo
}

// Login implements AuthService
func (authServ authService) CheckCredentials(creds model.CredentialsDto) error {
	savedCreds, err := authServ.userRepo.GetOne(creds.Username)
	if err == db.DatabaseError {
		return DatabaseError
	} else if err != nil {
		return err
	}

	return bcrypt.CompareHashAndPassword(savedCreds.PasswordHash, []byte(creds.Password))
}

// SaveCredentials implements AuthService
