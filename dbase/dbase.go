package dbase

import (
	_ "database/sql"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/joho/godotenv"
	"log"
	"os"
	"math/rand"
)

/*
	
*/

func InitDB() {
	db := OpenDB()
	defer db.Close()
	db.DB()
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Item{})
	db.AutoMigrate(&Order{})
}

func AddTestCode(db *gorm.DB) {
	AddUsers(db)
	AddItems(db)
}

func AddUsers(db *gorm.DB) {
	for i := 0; i < 10; i++ {
		user := User{CustomerCardID: fmt.Sprintf("cusCardID_%d", i+1), Name: fmt.Sprintf("userName_%d", i+1), Items: nil, Orders: nil}
		db.Create(&user)
	}
}

func AddItems(db *gorm.DB) {
	for i := 0; i < 10; i++ {
		item := Item{UserID: rand.Intn(4)+1, Name: fmt.Sprintf("itemName_%d", i+1), Orders: nil}
		db.Create(&item)
	}
}

func OpenDB() *gorm.DB {
	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	sql := fmt.Sprintf("%v@%v/%v?charset=utf8&parseTime=True", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))
	db, err := gorm.Open(os.Getenv("DB"), sql)
	if err != nil {
		fmt.Printf("エラー%v\n", err)
		return nil
	}
	return db
}
