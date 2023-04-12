package db

import (
	"apr-backend/internal/model"
	"database/sql"
	"errors"
	"fmt"
	"log"
)

var InvalidFilter = errors.New("Invalid filter")

func NewCompanyRepository(db *sql.DB) CompanyRepository {
	return companyRepository{
		db: db,
	}
}

type CompanyRepository interface {
	SaveCompany(com model.Company) error
	FindCompanies(filter model.CompanyFilter) ([]model.Company, error)
}

type companyRepository struct {
	db *sql.DB
}

func validateColumn(col string) bool {
	validColumns := []string{"naziv", "vlasnik", "PIB", "maticniBroj", "mesto", "sediste"}
	for _, valCol := range validColumns {
		if col == valCol {
			return true
		}
	}
	return false
}

// FindCompanies implements CompanyRepository
func (cr companyRepository) FindCompanies(filter model.CompanyFilter) ([]model.Company, error) {
	query := `SELECT PIB, delatnost, vlasnik, c.naziv, adresaSedista, postanskiBroj, mesto, maticniBroj, n.oznaka, n.naziv as nstjNaziv
        FROM company c
        LEFT JOIN NSTJ n ON c.sediste = n.oznaka 
        WHERE (? = "" OR delatnost = ?)
        AND (? = "" OR sediste = ?)
        AND (? = "" OR mesto = ?)
        ORDER BY `
	valid := validateColumn(filter.OrderBy)
	if !valid {
		return []model.Company{}, fmt.Errorf("%w: %s is an invalid column", InvalidFilter, filter.OrderBy)
	}

	if filter.Asc {
		query = fmt.Sprintf("%s %s %s", query, filter.OrderBy, "ASC")
	} else {
		query = fmt.Sprintf("%s %s %s", query, filter.OrderBy, "DESC")
	}

	query = fmt.Sprintf("%s %s %d;", query, "LIMIT 50 OFFSET", filter.Page*50)

	stmt, err := cr.db.Prepare(query)
	if err != nil {
		return []model.Company{}, DatabaseError
	}

	rows, err := stmt.Query(filter.Delatnost, filter.Delatnost, filter.Sediste, filter.Sediste, filter.Mesto, filter.Mesto)

	companies := make([]model.Company, 0, 50)
	for rows.Next() {
		var company model.Company
		err := rows.Scan(&company.PIB, &company.Delatnost, &company.Vlasnik, &company.Naziv, &company.AdresaSedista, &company.PostanskiBroj, &company.Mesto, &company.MaticniBroj, &company.Sediste.Oznaka, &company.Sediste.Naziv)
		if err != nil {
			return companies, fmt.Errorf("%w: couldn't scan company %#v", DatabaseError, company)
		}
		companies = append(companies, company)
	}
	return companies, nil
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
