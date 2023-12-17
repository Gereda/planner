package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "mysecretpassword"
	dbname   = "postgres"
)

func InitDB() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
		CREATE TABLE if not exists planner (
		ID SERIAl PRIMARY KEY,
		Title VARCHAR(255),
		Description VARCHAR(255),
		Status BOOLEAN,
		Priority VARCHAR(255)
		);
	`)
	if err != nil {
		fmt.Print(err)
		return nil
	}

	Prepare(db)

	return db
}

var (
	Insert    *sql.Stmt
	SelectAll *sql.Stmt
	//Select *sql.Stmt
	//PrGetTaskByID  *sql.Stmt
	//PrUpdateTasks  *sql.Stmt
	//PrDelTask      *sql.Stmt
)

func Prepare(db *sql.DB) (err error) {
	Insert, err = db.Prepare("INSERT INTO planner (title, description, status, priority) VALUES ($1, $2, $3, $4) RETURNING ID")
	if err != nil {
		return
	}

	SelectAll, err = db.Prepare("SELECT * FROM planner")
	if err != nil {
		return
	}

	//Select, err = db.Prepare("SELECT * FROM planner WHERE status = $1")
	//if err != nil {
	//	return
	//}
	//
	//PrGetTaskByID, err = db.Prepare("SELECT * FROM planner WHERE ID = $1")
	//if err != nil {
	//	return
	//}
	//
	//PrUpdateTasks, err = db.Prepare("UPDATE planner SET title = $1, description = $2, status = $3, priority = $4 WHERE ID = $5")
	//if err != nil {
	//	return
	//}
	//
	//PrDelTask, err = db.Prepare("DELETE FROM planner WHERE ID = $1")
	//if err != nil {
	//	return
	//}
	return
}
