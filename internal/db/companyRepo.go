package db

import (
	"apr-backend/internal/model"
	"database/sql"
	"errors"
	"fmt"
	"log"
)

var InvalidFilter = errors.New("Invalid filter")
var NoSuchPibError = errors.New("PIB not found in database")

func NewCompanyRepository(db *sql.DB, personRepo PersonRepository) CompanyRepository {
	return companyRepository{
		db:         db,
		personRepo: personRepo,
	}
}

type CompanyRepository interface {
	// Saves a company, making sure that person with JMBG already exists in the db
	SaveCompany(com *model.Company) error
	FindCompanies(filter model.CompanyFilter) ([]model.Company, error)
	FindOne(pib int) (model.Company, error)
	FindOneCredentials(pib int) (model.Company, error)
}

type companyRepository struct {
	db         *sql.DB
	personRepo PersonRepository
}

// FindOne implements CompanyRepository
func (cr companyRepository) FindOne(pib int) (model.Company, error) {
	query := `SELECT PIB, delatnost, vlasnik, c.naziv, adresaSedista,
    postanskiBroj, mesto, n.oznaka, n.naziv as nstjNaziv
    FROM company c
    LEFT JOIN NSTJ n ON c.sediste = n.oznaka
    WHERE c.PIB = ?`

	stmt, err := cr.db.Prepare(query)
	if err != nil {
		log.Printf("err: %s\n", err.Error())
		return model.Company{}, fmt.Errorf("Error creating prepared statement: %w", DatabaseError)
	}

	var company model.Company
	err = stmt.QueryRow(pib).Scan(&company.PIB, &company.Delatnost, &company.Vlasnik, &company.Naziv, &company.AdresaSedista, &company.PostanskiBroj, &company.Mesto, &company.Sediste.Oznaka, &company.Sediste.Naziv)
	if err == sql.ErrNoRows {
		return model.Company{}, fmt.Errorf("Company with PIB %d not found: %w", pib, NoSuchPibError)
	}
	if err != nil {
		log.Printf("Error getting company with pib %d: %s", pib, err.Error())
		return model.Company{}, DatabaseError
	}
	return company, nil
}

// FindOne implements CompanyRepository
func (cr companyRepository) FindOneCredentials(pib int) (model.Company, error) {
	query := `SELECT c.password
    FROM company c
    WHERE c.PIB = ?`

	stmt, err := cr.db.Prepare(query)
	if err != nil {
		log.Printf("Error creating prepared statement: %s\n", err.Error())
		return model.Company{}, fmt.Errorf("Error creating prepared statement: %w", DatabaseError)
	}

	var company model.Company
	err = stmt.QueryRow(pib).Scan(&company.Password)
	if err == sql.ErrNoRows {
		return model.Company{}, fmt.Errorf("Company with PIB %d not found: %w", pib, NoSuchPibError)
	}
	if err != nil {
		log.Printf("Error getting company with pib %d: %s", pib, err.Error())
		return model.Company{}, DatabaseError
	}
	return company, nil
}

func validateColumn(col string) bool {
	validColumns := []string{"naziv", "vlasnik", "PIB", "mesto"}
	for _, valCol := range validColumns {
		if col == valCol {
			return true
		}
	}
	return false
}

// FindCompanies implements CompanyRepository
func (cr companyRepository) FindCompanies(filter model.CompanyFilter) ([]model.Company, error) {
	query := `SELECT PIB, delatnost, vlasnik, c.naziv, adresaSedista, postanskiBroj, mesto, n.oznaka, n.naziv as nstjNaziv
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
		err := rows.Scan(&company.PIB, &company.Delatnost, &company.Vlasnik, &company.Naziv, &company.AdresaSedista, &company.PostanskiBroj, &company.Mesto, &company.Sediste.Oznaka, &company.Sediste.Naziv)
		if err != nil {
			return companies, fmt.Errorf("%w: couldn't scan company %#v", DatabaseError, company)
		}
		companies = append(companies, company)
	}
	return companies, nil
}

// SaveCompany implements CompanyRepository
func (cr companyRepository) SaveCompany(com *model.Company) error {
	tx, err := cr.db.Begin()
	if err != nil {
		return DatabaseError
	}
	defer tx.Rollback()

	_, err = cr.personRepo.GetOne(com.Vlasnik, tx)
	if err != nil {
		return fmt.Errorf("Error getting user with JMBG %s: %w", com.Vlasnik, NoSuchJmbgError)
	}

	stmt, err := tx.Prepare(`INSERT INTO company
        (delatnost, vlasnik, naziv, adresaSedista, postanskiBroj, mesto, sediste, password)
        VALUES(?, ?, ?, ?, ?, ?, ?, ?);`)
	if err != nil {
		log.Printf("Error when creating prepared statement: %s", err.Error())
		return fmt.Errorf("%w", DatabaseError)
	}
	defer stmt.Close()

	res, err := stmt.Exec(com.Delatnost, com.Vlasnik, com.Naziv, com.AdresaSedista, com.PostanskiBroj, com.Mesto, com.Sediste.Oznaka, com.Password)
	if err != nil {
		log.Printf("Insert error: %s", err.Error())
		return fmt.Errorf("%w", DatabaseError)
	}
	pib, err := res.LastInsertId()

	com.PIB = int(pib)

	if err != nil {
		return fmt.Errorf("Error when getting PIB of new company: %w", DatabaseError)
	}

	tx.Commit()
	return nil
}
