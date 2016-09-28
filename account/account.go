package account

import (
    "dbase"
    "math/rand"
    "github.com/jinzhu/gorm"
    "fmt"
)

func InitDB() {
	db := dbase.OpenDB()
	defer db.Close()
	db.DB()
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Item{})
	db.AutoMigrate(&Order{})
}

func AddTestCode(db *gorm.DB) {
	addTestUsers(db)
	addTestItems(db)
}
func addTestUsers(db *gorm.DB) {
	for i := 0; i < 10; i++ {
		user := User{CustomerCardID: fmt.Sprintf("cusCardID_%d", i+1), Name: fmt.Sprintf("userName_%d", i+1), Items: nil, Orders: nil}
		db.Create(&user)
	}
}
func addTestItems(db *gorm.DB) {
	for i := 0; i < 10; i++ {
		item := Item{
			UserID: uint(rand.Intn(4) + 1), 
			Name: fmt.Sprintf("itemName_%d", i+1), 
			BasePrice: i*1000, 
			DailyCharge: i*800,
			DepositFee: i*20000,
			Orders: nil}
		db.Create(&item)
	}
}