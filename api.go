package main

import (
	"Rest_PostgreSQL/toy"
	"flag"
	"fmt"
	//	"github.com/gorilla/mux"
	"gopkg.in/pg.v5"
	"log"
	"net/http"
	"os"
	//"path"
	"time"
)

//Enter point, read configuration and route all request
func main() {
	optionInitDB := flag.Bool("initDB", false, "Clean all data and initial all tables in DB")
	optionHost := flag.String("host", "192.168.100.1", "Postgre host addr")
	optionDbName := flag.String("db", "", "DB Name")
	optionUser := flag.String("user", "", "DB user name")
	optionDbPwd := flag.String("pwd", "", "Db Password")
	optionPort := flag.Int("port", 8080, "Port to listen on")
	optionLog := flag.String("log", "test.log", "name of logfile (if not supplied default to stdout)")
	//optionData := flag.String("data", "/tmpData", "location of temp directory for caching")
	flag.Parse()
	//if user forget set Db name and pwd, show a notification and exit.
	if (*optionDbName == "") || (*optionDbPwd == "") {
		fmt.Printf("Must specify Db credentials using: -db {DB Name} -pwd {password}\n")
		os.Exit(1)
	}

	//start goroutine to switch log file every day begin!
	if *optionLog != "" {
		go toy.LogRunner(*optionLog, 24*60*time.Minute)
	}

	//open db
	db := pg.Connect(&pg.Options{
		Addr:     *optionHost,
		User:     *optionUser,
		Database: *optionDbName,
		Password: *optionDbPwd,
	})
	//Init Db
	if *optionInitDB == true {
		err := CreateTable()
		if err != nil {
			log.Fatal(err)
		}
	}

	r := toy.RegisteRouters(db)
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *optionPort), nil))
	err := db.Close()
	if err != nil {
		panic(err)
	}
}
func CreateTable() error {
	return nil
}
