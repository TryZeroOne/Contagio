package database

import (
	"database/sql"
	"os"

	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

var Db *sql.DB

func InitDb(createUser bool) {
	db, err := sql.Open("sqlite3", "./sqlite/database.db")

	if err != nil {
		os.Create("./sqlite/database.db")
		InitDb(createUser)
	}

	Db = db

	statement, err := db.Prepare("CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, login TEXT, password TEXT)")

	if err != nil {
		fmt.Println(err)
		return
	}

	statement.Exec()

	statement, err = db.Prepare("CREATE TABLE IF NOT EXISTS allowed (id INTEGER PRIMARY KEY, ip TEXT)")

	if err != nil {
		fmt.Println(err)
		return
	}

	statement.Exec()

	if createUser {

		statement, err = db.Prepare("INSERT INTO users(login, password) VALUES (?, ?)")

		if err != nil {
			fmt.Println(err)
			return
		}

		statement.Exec("root", "root")
	}
}
