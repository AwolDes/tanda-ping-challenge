package pings

import (
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
	"database/sql"
	"github.com/go-sql-driver/mysql"
)


func CreateDevice(w http.ResponseWriter, r *http.Request) {}

func GetDeviceOnDate(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "user:password@/dbname")
	devices := make([]Device, 0)
	device := &Device {
		Id: "abc123",
		Timestamps: []int{123, 1234, 111},
	}
	devices = append(devices, *device)

	params := mux.Vars(r)
	for _, item := range devices {
		if item.Id == params["deviceId"] {
			json.NewEncoder(w).Encode(item)
		}
	}
}

func GetDeviceDateRange(w http.ResponseWriter, r *http.Request) {}
func GetAllDevicesOnDate(w http.ResponseWriter, r *http.Request) {}
func GetAllDevicesInDateRange(w http.ResponseWriter, r *http.Request) {}

func GetAllDevices(w http.ResponseWriter, r *http.Request) {
	devices := make([]Device, 0)
	device := &Device {
		Id: "abc123",
		Timestamps: []int{123, 1234, 111},
	}
	devices = append(devices, *device)
	json.NewEncoder(w).Encode(devices)
}

func ClearData(w http.ResponseWriter, r *http.Request) {
	// delete all data from table
	res := &HttpRes{
		StatusCode: 200,
	}
	json.NewEncoder(w).Encode(res)
}