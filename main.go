package main

import (
    _ "account"
    "fmt"
    "github.com/joho/godotenv"
    "log"
    "os"
    "dbase"
)

func Env_load() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
    Env_load()
    dbase.InitDB()
    db := dbase.OpenDB()
    dbase.AddTestCode(db)
    //account.TestDB()
    fmt.Println(os.Getenv("DB_USER"))
}