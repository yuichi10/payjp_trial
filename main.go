package main

import (
	"account"
	"fmt"
	"github.com/gorilla/context"
	_ "github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	_ "os"
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
	account.InitDB()
	
	//account.TestDB()
	db := dbase.OpenDB()
	defer db.Close()
	//account.AddTestCode(db)
	http.HandleFunc("/", helloGo)
	http.HandleFunc("/account/publishOrder", account.PublishOrder)
	http.HandleFunc("/account/getToken", account.GetToken)
	http.HandleFunc("/account/addCardInfo", account.AddCardInfo)
	http.HandleFunc("/testDB", account.TestDB)
	http.ListenAndServe(":9978", context.ClearHandler(http.DefaultServeMux))
}

func helloGo(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello GO!!!!")
}
