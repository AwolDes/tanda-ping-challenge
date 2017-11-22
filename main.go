package main

import (
    "log"
    "net/http"
	"github.com/gorilla/mux"
	"pings"
)

func main() {

	pings.DbConnect()

	print("Running server at http://localhost:8000\n")
	router := mux.NewRouter()
    router.HandleFunc("/{deviceId}/{timestamp}", pings.CreateDevicePing).Methods("POST")
	router.HandleFunc("/{deviceId}/{date}", pings.GetDeviceOnDate).Methods("GET")
	router.HandleFunc("/{deviceId}/{from}/{to}", pings.GetDeviceDateRange).Methods("GET")
    router.HandleFunc("/all/{date}", pings.GetAllDevicesOnDate).Methods("GET")
	router.HandleFunc("/all/{from}/{to}", pings.GetAllDevicesInDateRange).Methods("GET")
	router.HandleFunc("/devices", pings.GetAllDevices).Methods("GET")
	router.HandleFunc("/clear_data", pings.ClearData).Methods("POST")
    log.Fatal(http.ListenAndServe(":8000", router))
}