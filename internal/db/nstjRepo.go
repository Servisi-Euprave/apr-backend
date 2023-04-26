package db

import (
	"apr-backend/internal/model"
	"database/sql"
	"fmt"
	"log"
)

type NstjRepository interface {
	FindAll() ([]model.Nstj, error)
}

func NewNstjRepository(db *sql.DB) NstjRepository {
	return nstjRepo{
		db: db,
	}
}

type nstjRepo struct {
	db *sql.DB
}

// FindAll implements NstjRepository
func (nstjRepo nstjRepo) FindAll() ([]model.Nstj, error) {
	query := `SELECT oznaka, naziv FROM apr.NSTJ;`
	rows, err := nstjRepo.db.Query(query)
	if err != nil {
		log.Printf("Couldn't prepare statement: %s", err.Error())
		return []model.Nstj{}, fmt.Errorf("Couldn't prepare statement: %w", DatabaseError)
	}
	defer rows.Close()

	nstjCollection := make([]model.Nstj, 0, 100)
	for rows.Next() {
		var nstj model.Nstj
		rows.Scan(&nstj.Oznaka, &nstj.Naziv)
		nstjCollection = append(nstjCollection, nstj)
	}
	if rows.Err() != nil {
		log.Printf("Error reading NSTJ: %s\n", rows.Err().Error())
		err = fmt.Errorf("Error reading NSTJ: %w", DatabaseError)
	}
	return nstjCollection, err
}
