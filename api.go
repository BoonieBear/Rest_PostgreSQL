package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"gopkg.in/pg.v5"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"
)

//Enter point, read configuration and route all request
func main() {
	optionVerbose := flag.Bool("verb", false, "Verbose output for API handlers")
	optionVerbose := flag.Bool("initDB", false, "Clean all data and initial all tables in DB")
	optionHost := flag.String("host", "192.168.100.1", "Postgre host addr")
	optionDbName := flag.String("db", "", "DB Name")
	optionDbPwd := flag.String("pwd", "", "Db Password")
	optionPort := flag.Int("port", 8080, "Port to listen on")
	optionLog := flag.String("log", "test.log", "name of logfile (if not supplied default to stdout)")
	optionData := flag.String("data", "/tmpData", "location of temp directory for caching")
	flag.Parse()
	//if user forget set Db name and pwd, show a notification and exit.
	if (optionDbName == "") || (optionDbPwd == "") {
		fmt.Printf("Must specify Db credentials using: -db {DB Name} -pwd {password}\n")
		os.Exit(1)
	}
	verb := *optionVerbose

	//start goroutine to switch log file every day begin!
	if *optionLog != "" {
		go LogRunner(*optionLog, 24*60*time.Minute)
	}
	port := *optionPort

	//open db

}
