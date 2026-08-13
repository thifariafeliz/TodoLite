package main

import (
	"fmt"
	"errors"
	"database/sql"
	_ "modernc.org/sqlite"
)

type Domain struct {
	ID int
	Name string
}

type Task struct {
	ID int
	Title string
	Done bool
	DomainID int
}

func ConnectDB(dbname string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dbname)
	if err != nil {
		fmt.Println("Failed to open database:", err)
		return nil, err
	}

	if db == nil {
		fmt.Println("Database is nil")
		return nil, errors.New("Database is nil")
	}

	err = db.Ping()
	if err != nil {
		fmt.Println("Failed to ping database:", err)
		return nil, err
	}

	return db, nil
}

func CreateTables(db *sql.DB) error {
	if db == nil {
		fmt.Println("Bad argument to CreateTables: database is nil.")
		return errors.New("Database is nil")
	}

	err := db.Ping()
	if err != nil {
		fmt.Println("Failed to ping database:", err)
		return err
	}

	TaskTableSQL := `
		CREATE TABLE IF NOT EXISTS tasks (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			done BOOL NOT NULL
		);

		CREATE TABLE IF NOT EXISTS domains (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT UNIQUE NOT NULL
		);

		CREATE TABLE IF NOT EXISTS task_domain (
			task_id INTEGER,
			domain_id INTEGER,
			PRIMARY KEY (task_id, domain_id),
			FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE,
			FOREIGN KEY (domain_id) REFERENCES domains(id) ON DELETE CASCADE
		);
	`

	stmt, err := db.Prepare(TaskTableSQL)
	if err != nil {
		fmt.Println("Failed to prepare SQL tables creation statement:", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		fmt.Println("Failed to exec table creation SQL statement:", err)
		return err
	}

	return nil
}

func CheckTableExistence(tbname string) error {
	
}

func PrintTasks() {}

func AddTask() {}

func main() {
	db, err := ConnectDB("dbzinho.db")
	if err != nil {
		fmt.Println("It was not possible to connect to databse.")
		return
	}
	defer db.Close()

	fmt.Println("Connection with database made successfully")

	

}
