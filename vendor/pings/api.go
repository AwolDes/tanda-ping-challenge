package pings

import (
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
	"os"
	_ "github.com/joho/godotenv/autoload"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
)

var db *sql.DB 
var err error

func DbConnect() {
	DbPass := os.Getenv("DB_PASS")
	db, err = sql.Open("mysql", "root:" + DbPass + "@tcp(localhost:3306)/tanda_devices")
	if err != nil {
        panic(err.Error())    
	}
	
	// defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err.Error())
	} else {
		print("\nConnected To Database!\n")
	}
}

func checkErr(err error) {
    if err != nil {
        panic(err)
    }
}

func CreateDevicePing(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	deviceName := params["deviceId"]
	timestamp := params["timestamp"]
	timestampInt, _ := strconv.ParseInt(timestamp, 10, 64)
	stmt, err := db.Prepare("INSERT devices SET device_name=?,timestamp=?")
	checkErr(err)
	
	_, err = stmt.Exec(deviceName, timestampInt)
	
	if err != nil {
		res := &HttpRes{
			StatusCode: 400,
		}
		json.NewEncoder(w).Encode(res)
		checkErr(err)
		db.Close()
	} else {
		resSuccess := &HttpRes{
			StatusCode: 200,
		}
	
		json.NewEncoder(w).Encode(resSuccess)
		db.Close()
	}

}

func GetDeviceOnDate(w http.ResponseWriter, r *http.Request) {
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
	// devices := make([]Device, 0)
	// device := &Device {
	// 	Id: "abc123",
	// 	Timestamps: []int{123, 1234, 111},
	// }
	// devices = append(devices, *device)

	err := db.QueryRow("SELECT * FROM devices")

	if err != nil {
		res := &HttpRes{
			StatusCode: 400,
		}
		json.NewEncoder(w).Encode(res)
	}

	//json.NewEncoder(w).Encode(devices)
}

func ClearData(w http.ResponseWriter, r *http.Request) {
	// delete all data from table
	res := &HttpRes{
		StatusCode: 200,
	}
	json.NewEncoder(w).Encode(res)
}