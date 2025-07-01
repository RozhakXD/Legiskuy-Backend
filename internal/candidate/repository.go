package candidate

import (
	"database/sql"
	"legiskuy-backend/pkg/database"
)

type Candidate struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Party string `json:"party"`
	Votes int    `json:"votes"`
}

type Repository interface {
	Create(candidate *Candidate) (int64, error)
	FindAll(name, party string) ([]Candidate, error)
	FindByID(id int) (*Candidate, error)
	Update(id int, candidate *Candidate) error
	Delete(id int) error
	IncrementVoteCount(tx *sql.Tx, candidateID int) error
}

type repository struct {
	db *sql.DB
}

func NewRepository() Repository {
	return &repository{
		db: database.DB,
	}
}

func (r *repository) Create(candidate *Candidate) (int64, error) {
	query := `INSERT INTO candidates (name, party) VALUES (?, ?)`
	result, err := r.db.Exec(query, candidate.Name, candidate.Party)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *repository) FindAll(name, party string) ([]Candidate, error) {
	query := `SELECT id, name, party, votes FROM candidates WHERE 1=1`
	args := []interface{}{}

	if name != "" {
		query += " AND name LIKE ?"
		args = append(args, "%"+name+"%")
	}

	if party != "" {
		query += " AND party LIKE ?"
		args = append(args, "%"+party+"%")
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	candidates := make([]Candidate, 0)
	for rows.Next() {
		var c Candidate
		if err := rows.Scan(&c.ID, &c.Name, &c.Party, &c.Votes); err != nil {
			return nil, err
		}
		candidates = append(candidates, c)
	}
	return candidates, nil
}

func (r *repository) FindByID(id int) (*Candidate, error) {
	query := `SELECT id, name, party, votes FROM candidates WHERE id = ?`
	row := r.db.QueryRow(query, id)

	var c Candidate
	if err := row.Scan(&c.ID, &c.Name, &c.Party, &c.Votes); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil
}

func (r *repository) Update(id int, candidate *Candidate) error {
	query := `UPDATE candidates SET name = ?, party = ? WHERE id = ?`
	result, err := r.db.Exec(query, candidate.Name, candidate.Party, id)
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
	query := `DELETE FROM candidates WHERE id = ?`
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

func (r *repository) IncrementVoteCount(tx *sql.Tx, candidateID int) error {
	query := `UPDATE candidates SET votes = votes + 1 WHERE id = ?`
	_, err := tx.Exec(query, candidateID)
	return err
}
