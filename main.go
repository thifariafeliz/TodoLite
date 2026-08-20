package main

import (
	"fmt"
	"errors"
	"database/sql"
	_ "modernc.org/sqlite"
)

type Domain struct {
	Name string
	ID   int
}

type Task struct {
	Title    string
	DomainID []int
	ID       int
	Status   bool
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
			status BOOL NOT NULL
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

func CheckTableExistence(db *sql.DB, tbname string) (bool, error) {
	sqlTasks := `SELECT name FROM sqlite_master WHERE type='table' and name=?`

	var name string
	err := db.QueryRow(sqlTasks, tbname).Scan(&name)

	if err == sql.ErrNoRows {
		fmt.Println("Não existe nenhuma tabela chamada", tbname)
		return false, err
	}

	if err != nil {
		fmt.Println("Failed to check the table existence:", err)
		return false, err
	}

	return true, nil
}

func TakeTasksByStatus(db *sql.DB, done bool) ([]Task, error) {
	sql := ``

	if done {
		sql = `SELECT * FROM tasks WHERE done=true;`
	} else {
		sql = `SELECT * FROM tasks WHERE done=false;`
	}

	rows, err := db.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task

	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.ID, &task.Title, &task.DomainID, &task.Status); err != nil {
			return tasks, err
		}

		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		return tasks, err
	}

	return tasks, nil
}

// Adds a task to the database
func AddTask(db *sql.DB, task Task) (int64, error) {
	// adds a task to the database
	result, err := db.Exec("INSERT INTO tasks (title, status) VALUES (?, ?);", task.Title, task.Status)
	if err != nil {
		return 0, fmt.Errorf("AddTask: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("AddTask: %v", err)
	}

	return id, nil
}

// Adds domains to the domains' table
func AddDomain(db *sql.DB, domains []Domain) (int64, error) {
	domainNo := 0
	var id int64
	for i := 0; i < len(domains); i++ {
		// Checks if the domain beeing iterated is already in the table
		if err := db.QueryRow("SELECT COUNT(id) FROM domains WHERE name = '?'", domains[i].Name).Scan(&domainNo); err != nil {
			// If there is a row, then continue
			if err != sql.ErrNoRows {
				continue
			}
			return 0, err
		}

		if domainNo == 0 {
			result, err := db.Exec("INSERT INTO domains (name) VALUES (?);", domains[i].Name)
			if err != nil {
				return 0, err
			}

			id, err = result.LastInsertId()
			if err != nil {
				return 0, err
			}

			return id, nil
		}
	}

	return id, nil
}

func AddTaskDomain(db *sql.DB, taskID, domainID int) (int64, error) {
	// Need to add the taskID and domainID to the table of task_domain
	// 

	sqlExists := `SELECT COUNT(task_id) FROM task_domain WHERE task_id=? AND domain_id=?`
	
	exists := 0
	if err := db.QueryRow(sqlExists, taskID, domainID).Scan(&exists); err != nil {
		
	}
}

func main() {
	db, err := ConnectDB("dbzinho.db")
	if err != nil {
		fmt.Println("It was not possible to connect to databse.")
		return
	}
	defer db.Close()

	fmt.Println("Connection with database made successfully")

	err = CreateTables(db)
	if err != nil {
		fmt.Println("Failed to create tables:", err)
		return
	}

	exist, err := CheckTableExistence(db, "tasks")
	if exist == true {
		fmt.Println("Tasks table exists")
	}
	exist, err = CheckTableExistence(db, "domains")
	if exist == true {
		fmt.Println("Domains table exists")
	}

	exist, err = CheckTableExistence(db, "task_domain")
	if exist == true {
		fmt.Println("Task_Domain table exist")
	}

	task := Task{Title: "Limpar machucado.", Status: false}

}
