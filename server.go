package main

import (
	"database/sql"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/notlcry/galera-healthcheck/healthcheck"
	. "github.com/notlcry/galera-healthcheck/logger"
)

var serverPort = flag.Int(
	"port",
	8080,
	"Specifies the port of the healthcheck server",
)

var mysqlHost = flag.String(
	"host",
	"127.0.0.1",
	"Specifies the MySQL host to connect as",
)

var mysqlPort = flag.String(
	"mport",
	"3306",
	"Specifies the MySQL port to connect as",
)

var mysqlUser = flag.String(
	"user",
	"root",
	"Specifies the MySQL user to connect as",
)

var mysqlPassword = flag.String(
	"password",
	"",
	"Specifies the MySQL password to connect with",
)

var availableWhenDonor = flag.Bool(
	"availWhenDonor",
	true,
	"Specifies if the healthcheck allows availability when in donor state",
)

var availableWhenReadOnly = flag.Bool(
	"availWhenReadOnly",
	false,
	"Specifies if the healthcheck allows availability when in read only mode",
)

var pidfile = flag.String(
	"pidfile",
	"",
	"Location for the pidfile",
)

var connectionCutterPath = flag.String(
	"connectionCutterPath",
	"",
	"Location for the script which cuts mysql connections",
)

var healthchecker *healthcheck.Healthchecker

func handler(w http.ResponseWriter, r *http.Request) {
	result, msg := healthchecker.Check()
	if result != nil && result.Healthy {
		w.WriteHeader(http.StatusOK)
	} else if result != nil && !result.Healthy {
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		w.WriteHeader(http.StatusContinue)
	}

	fmt.Fprintf(w, "Galera Cluster Node status: %s", msg)
	LogWithTimestamp(msg)
}

func handlerSeq(w http.ResponseWriter, r *http.Request) {
	result, msg := healthchecker.Check()
	if result != nil && result.Healthy {
		w.WriteHeader(http.StatusOK)
	} else if result != nil && !result.Healthy {
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		w.WriteHeader(http.StatusContinue)
	}

	fmt.Fprintf(w, "seq is: %s", msg)
	LogWithTimestamp(msg)
}

func main() {
	flag.Parse()

	err := ioutil.WriteFile(*pidfile, []byte(strconv.Itoa(os.Getpid())), 0644)
	if err != nil {
		panic(err)
	}

	db, _ := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/", *mysqlUser, *mysqlPassword, *mysqlHost, *mysqlPort))
	config := healthcheck.HealthcheckerConfig{
		*availableWhenDonor,
		*availableWhenReadOnly,
	}

	healthchecker = healthcheck.New(db, config)

	http.HandleFunc("/", handler)
	http.HandleFunc("/seq", handlerSeq)
	http.ListenAndServe(fmt.Sprintf(":%d", *serverPort), nil)
}
