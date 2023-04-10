package db

import (
	"apr-backend/internal/model"
	"database/sql"
	"fmt"
	"log"
)

func NewCompanyRepository(db *sql.DB) CompanyRepository {
	return companyRepository{
		db: db,
	}
}

type CompanyRepository interface {
	SaveCompany(com model.Company) error
}

type companyRepository struct {
	db *sql.DB
}

// SaveCompany implements CompanyRepository
func (cr companyRepository) SaveCompany(com model.Company) error {
	stmt, err := cr.db.Prepare(`INSERT INTO company
        (PIB, delatnost, vlasnik, naziv, adresaSedista, postanskiBroj, mesto, sediste, maticniBroj)
        VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?);`)
	if err != nil {
		log.Printf("Error when creating prepared statement: %s", err.Error())
		return fmt.Errorf("%w", DatabaseError)
	}
	defer stmt.Close()

	_, err = stmt.Exec(com.PIB, com.Delatnost, com.Vlasnik, com.Naziv, com.AdresaSedista, com.PostanskiBroj, com.Mesto, com.Sediste.Oznaka, com.MaticniBroj)
	if err != nil {
		log.Printf("Insert error: %s", err.Error())
		return fmt.Errorf("%w", DatabaseError)
	}
	return nil
}
