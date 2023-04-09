package db

import (
	"apr-backend/internal/model"
	"database/sql"
	"errors"
	"fmt"
	"log"
)

func NewUserRepo(db *sql.DB) UserRepo {
	return userRepo{db: db}
}

type UserRepo interface {
	GetOne(username string) (model.Credentials, error)
	SaveUser(user model.User) error
}

type userRepo struct {
	db *sql.DB
}

var DatabaseError = errors.New("Error has occured when connecting to database")
var ErrUsernameTaken = errors.New("Username taken")

// Login implements AuthDB
func (userRepo userRepo) GetOne(username string) (model.Credentials, error) {
	creds := model.Credentials{}
	stmt, err := userRepo.db.Prepare("SELECT username, password FROM users WHERE username = ?")
	if err != nil {
		log.Printf("Error: %s", err.Error())
		return creds, DatabaseError
	}
	defer stmt.Close()

	err = stmt.QueryRow(username).Scan(&creds.Username, &creds.PasswordHash)
	if err == sql.ErrNoRows {
		log.Printf("Error: %s", err.Error())
		return creds, fmt.Errorf("User with the username %s does not exist", username)
	}
	if err != nil {
		log.Printf("Error: %s", err.Error())
		return creds, DatabaseError
	}

	return creds, nil
}

// SaveCredentials implements AuthDB
func (userRepo userRepo) SaveUser(user model.User) error {
	stmt, err := userRepo.db.Prepare(`INSERT INTO users (phone, email, sex, address, name, lastname, password, username, jmbg) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`)

	if err != nil {
		log.Printf("Error: %s", err.Error())
		return DatabaseError
	}
	defer stmt.Close()

	qRes, err := stmt.Exec(user.Phone, user.Email, user.Sex, user.Address, user.Name, user.Lastname, user.Password, []byte(user.Username), user.Jmbg)
	if err != nil {
		log.Printf("Error: %s", err.Error())
		return DatabaseError
	}

	c, err := qRes.RowsAffected()
	if err != nil {
		log.Printf("Error: %s", err.Error())
		return DatabaseError
	}

	if c == 0 {
		return ErrUsernameTaken
	}
	return nil
}
