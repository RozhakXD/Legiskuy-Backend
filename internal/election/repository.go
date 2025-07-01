package election

import (
	"database/sql"
	"legiskuy-backend/pkg/database"
)

type Repository interface {
	BeginTransaction() (*sql.Tx, error)
	CreateVote(tx *sql.Tx, voterID, candidateID int) error

	GetSetting(key string) (string, error)
	SetSetting(key, value string) error
}

type repository struct {
	db *sql.DB
}

func NewRepository() Repository {
	return &repository{
		db: database.DB,
	}
}

func (r *repository) BeginTransaction() (*sql.Tx, error) {
	return r.db.Begin()
}

func (r *repository) CreateVote(tx *sql.Tx, voterID, candidateID int) error {
	query := `INSERT INTO votes (voter_id, candidate_id) VALUES (?, ?)`
	_, err := tx.Exec(query, voterID, candidateID)
	return err
}

func (r *repository) GetSetting(key string) (string, error) {
	query := `SELECT value FROM settings WHERE key = ?`
	row := r.db.QueryRow(query, key)
	var value string
	err := row.Scan(&value)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", err
	}
	return value, nil
}

func (r *repository) SetSetting(key, value string) error {
	queryUpdate := `UPDATE settings SET value = ? WHERE key = ?`
	res, err := r.db.Exec(queryUpdate, value, key)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		queryInsert := `INSERT INTO settings (key, value) VALUES (?, ?)`
		_, err := r.db.Exec(queryInsert, key, value)
		return err
	}
	return nil
}
