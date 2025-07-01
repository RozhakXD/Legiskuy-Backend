package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func ConnectDB() {
	var err error

	DB, err = sql.Open("sqlite3", "./legiskuy.db")
	if err != nil {
		log.Fatal("Gagal terhubung ke database:", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("Gagal ping ke database:", err)
	}

	log.Println("Berhasil terhubung ke database.")

	createTables()
}

func createTables() {
	candidatesTable := `
	CREATE TABLE IF NOT EXISTS candidates (
		"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"name" TEXT NOT NULL,
		"party" TEXT NOT NULL,
		"votes" INTEGER DEFAULT 0
	);`

	votersTable := `
    CREATE TABLE IF NOT EXISTS voters (
        "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
        "name" TEXT NOT NULL,
        "has_voted" BOOLEAN DEFAULT FALSE
    );`

	votesTable := `
	CREATE TABLE IF NOT EXISTS votes (
		"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"voter_id" INTEGER,
		"candidate_id" INTEGER,
		"created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(voter_id) REFERENCES voters(id),
		FOREIGN KEY(candidate_id) REFERENCES candidates(id)
	);`

	settingsTable := `
	CREATE TABLE IF NOT EXISTS settings (
		"key" TEXT NOT NULL PRIMARY KEY,
		"value" TEXT
	);`

	usersTable := `
    CREATE TABLE IF NOT EXISTS users (
        "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
        "name" TEXT NOT NULL,
        "username" TEXT NOT NULL UNIQUE,
        "password" TEXT NOT NULL,
        "role" TEXT NOT NULL,
        "has_voted" BOOLEAN DEFAULT FALSE
    );`

	if _, err := DB.Exec(candidatesTable); err != nil {
		log.Fatal("Gagal membuat tabel candidates:", err)
	}
	if _, err := DB.Exec(votersTable); err != nil {
		log.Fatal("Gagal membuat tabel voters:", err)
	}
	if _, err := DB.Exec(usersTable); err != nil {
		log.Fatal("Gagal membuat tabel users:", err)
	}
	if _, err := DB.Exec(votesTable); err != nil {
		log.Fatal("Gagal membuat tabel votes:", err)
	}
	if _, err := DB.Exec(settingsTable); err != nil {
		log.Fatal("Gagal membuat tabel settings:", err)
	}

	log.Println("Tabel berhasil dibuat atau sudah ada.")
}
