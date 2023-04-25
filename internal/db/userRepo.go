package db

import (
	"apr-backend/internal/model"
	"database/sql"
	"errors"
	"fmt"
	"log"
)

func NewUserRepo(db *sql.DB) PersonRepository {
	return personRepo{db: db}
}

type PersonRepository interface {
	GetOne(jmbg string, tx *sql.Tx) (model.Person, error)
}

type personRepo struct {
	db *sql.DB
}

var DatabaseError = errors.New("Error has occured when connecting to database")
var NoSuchJmbgError = errors.New("JMBG doesn't exist in the database")

// Login implements AuthDB
func (pr personRepo) GetOne(jmbg string, tx *sql.Tx) (model.Person, error) {
	person := model.Person{}
	stmt, err := tx.Prepare("SELECT jmbg, name, lastname FROM person WHERE jmbg = ?")
	if err != nil {
		log.Printf("Error: %s", err.Error())
		return person, DatabaseError
	}
	defer stmt.Close()

	err = stmt.QueryRow(jmbg).Scan(&person.Jmbg, &person.Name, &person.Lastname)
	if err == sql.ErrNoRows {
		log.Printf("Error: %s", err.Error())
		return person, fmt.Errorf("Person with jmbg %s does not exist: %w", jmbg, NoSuchJmbgError)
	}
	if err != nil {
		log.Printf("Error: %s", err.Error())
		return person, DatabaseError
	}

	return person, nil
}
