package services

import (
	"apr-backend/internal/db"
	"apr-backend/internal/model"

	"golang.org/x/crypto/bcrypt"
)

func NewUserService(userRepo db.UserRepo) UserService {
	return userService{userRepo: userRepo}
}

type UserService interface {
	SaveUser(user model.User) error
}

type userService struct {
	userRepo db.UserRepo
}

// SaveUser implements UserService
func (userServ userService) SaveUser(user model.User) error {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcryptCost)
	if err != nil {
		return err
	}
	user.Password = hashedPass
	if err = userServ.userRepo.SaveUser(user); err != nil {
		if err == db.DatabaseError {
			return DatabaseError
		}
		return err
	}
	return nil
}
