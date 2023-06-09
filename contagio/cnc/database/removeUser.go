package database

import (
	"contagio/contagio/cnc/utils"
	"fmt"
)

func RemoveUser(login string) {
	statement, err := Db.Prepare("DELETE from users where login=?")
	if err != nil {
		fmt.Println(err)
	}
	statement.Exec(utils.Sha3(login))

}
