package account

import (
	"dbase"
	"fmt"
)


func TestDB() {
	db := dbase.OpenDB()
	defer db.Close()
	db.DB()
	orderType := new(dbase.Order)
	db.First(orderType)
	fmt.Println(orderType)
}
