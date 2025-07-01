package voter

import (
	"database/sql"
	"legiskuy-backend/pkg/database"
)

type Voter struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	HasVoted bool   `json:"has_voted"`
}

type Repository interface {
	Create(voter *Voter) (int64, error)
	FindAll(name string) ([]Voter, error)
	FindByID(id int) (*Voter, error)
	Update(id int, voter *Voter) error
	Delete(id int) error
	MarkAsVoted(tx *sql.Tx, VoterID int) error
}

type repository struct {
	db *sql.DB
}

func NewRepository() Repository {
	return &repository{
		db: database.DB,
	}
}

func (r *repository) Create(voter *Voter) (int64, error) {
	query := `INSERT INTO voters (name) VALUES (?)`
	result, err := r.db.Exec(query, voter.Name)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *repository) FindAll(name string) ([]Voter, error) {
	query := `SELECT id, name, has_voted FROM voters WHERE 1=1`
	args := []interface{}{}

	if name != "" {
		query += " AND name LIKE ?"
		args = append(args, "%"+name+"%")
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var voters []Voter
	for rows.Next() {
		var v Voter
		if err := rows.Scan(&v.ID, &v.Name, &v.HasVoted); err != nil {
			return nil, err
		}
		voters = append(voters, v)
	}
	return voters, nil
}

func (r *repository) FindByID(id int) (*Voter, error) {
	query := `SELECT id, name, has_voted FROM voters WHERE id = ?`
	row := r.db.QueryRow(query, id)

	var v Voter
	if err := row.Scan(&v.ID, &v.Name, &v.HasVoted); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &v, nil
}

func (r *repository) Update(id int, voter *Voter) error {
	query := `UPDATE voters SET name = ? WHERE id = ?`
	result, err := r.db.Exec(query, voter.Name, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *repository) Delete(id int) error {
	query := `DELETE FROM voters WHERE id = ?`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *repository) MarkAsVoted(tx *sql.Tx, voterID int) error {
	query := `UPDATE voters SET has_voted = TRUE WHERE id = ?`
	_, err := tx.Exec(query, voterID)
	return err
}
