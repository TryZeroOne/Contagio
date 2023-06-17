package database

import "fmt"

func CheckUser(login, password string) bool {
	rows, err := Db.Query("SELECT id, login, password FROM users WHERE login=? AND password=?", login, password)

	if err != nil {
		fmt.Println("[contagio] Database: " + err.Error())
		return false
	}

	defer rows.Close()
	var id int
	var _login string
	var _password string

	for rows.Next() {
		rows.Scan(&id, &_login, &_password)
		return true
	}
	return false
}
