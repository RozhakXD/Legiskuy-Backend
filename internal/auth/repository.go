package auth

import (
	"database/sql"
	"legiskuy-backend/pkg/database"
)

type Repository interface {
	Create(user *User) (*User, error)
	FindByUsername(username string) (*User, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository() Repository {
	return &repository{
		db: database.DB,
	}
}

func (r *repository) Create(user *User) (*User, error) {
	query := `INSERT INTO users (name, username, password, role) VALUES (?, ?, ?, ?)`
	result, err := r.db.Exec(query, user.Name, user.Username, user.Password, user.Role)
	if err != nil {
		return nil, err
	}
	id, _ := result.LastInsertId()
	user.ID = int(id)
	return user, nil
}

func (r *repository) FindByUsername(username string) (*User, error) {
	query := `SELECT id, name, username, password, role, has_voted FROM users WHERE username = ?`
	row := r.db.QueryRow(query, username)
	var u User
	err := row.Scan(&u.ID, &u.Name, &u.Username, &u.Password, &u.Role, &u.HasVoted)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}
