package main

import (
    "log"
    "net/http"
	"github.com/gorilla/mux"
	"pings"
	"os"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func DbConnect()  *sql.DB {
	DbPass := os.Getenv("DB_PASS")
	db, err := sql.Open("mysql", "root:" + DbPass + "@tcp(localhost:3306)/tanda_devices")
	if err != nil {
        panic(err.Error())    
	}

	err = db.Ping()
	if err != nil {
		panic(err.Error())
	} else {
		print("\nConnected To Database!\n")
		return db
	}
}

func main() {
	api := &pings.Api{Db: DbConnect()}
	db := api.Db

	// Long lived connection
	// Go prefers global variables
	defer db.Close()
	

	print("Running server at http://localhost:8000\n")
	router := mux.NewRouter()
    router.HandleFunc("/all/{date}", api.GetAllDevicesOnDate).Methods("GET")
	router.HandleFunc("/all/{from}/{to}", api.GetAllDevicesInDateRange).Methods("GET")
    router.HandleFunc("/{deviceId}/{timestamp}", api.CreateDevicePing).Methods("POST")
	router.HandleFunc("/{deviceId}/{date}", api.GetDeviceOnDate).Methods("GET")
	router.HandleFunc("/{deviceId}/{from}/{to}", api.GetDeviceDateRange).Methods("GET")
	router.HandleFunc("/devices", api.GetAllDevices).Methods("GET")
	router.HandleFunc("/clear_data", api.ClearData).Methods("POST")
    log.Fatal(http.ListenAndServe(":8000", router))
}