package pings

import (
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
)

func checkErr(err error) {
    if err != nil {
        panic(err)
    }
}

func (api *Api) CreateDevicePing(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	deviceName := params["deviceId"]
	timestamp := params["timestamp"]
	timestampInt, _ := strconv.ParseInt(timestamp, 10, 64)
	stmt, err := api.Db.Prepare("INSERT devices SET device_name=?,timestamp=?")
	checkErr(err)
	
	_, err = stmt.Exec(deviceName, timestampInt)
	
	if err != nil {
		res := &HttpRes{
			StatusCode: 400,
		}
		json.NewEncoder(w).Encode(res)
		checkErr(err)
	} else {
		resSuccess := &HttpRes{
			StatusCode: 200,
		}
	
		json.NewEncoder(w).Encode(resSuccess)
	}

}

func (api *Api) GetDeviceOnDate(w http.ResponseWriter, r *http.Request) {
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

func (api *Api) GetDeviceDateRange(w http.ResponseWriter, r *http.Request) {}
func (api *Api) GetAllDevicesOnDate(w http.ResponseWriter, r *http.Request) {}
func (api *Api) GetAllDevicesInDateRange(w http.ResponseWriter, r *http.Request) {}

func (api *Api) GetAllDevices(w http.ResponseWriter, r *http.Request) {

	rows, err := api.Db.Query("SELECT DISTINCT device_name FROM devices")

	if err != nil {
		res := &HttpRes{
			StatusCode: 400,
		}

		json.NewEncoder(w).Encode(res)
		checkErr(err)
	}

	defer rows.Close()

	deviceNames := make([]string, 0)

	for rows.Next() {
		var deviceName string
		rows.Scan(&deviceName)
		deviceNames = append(deviceNames, deviceName)

	}

	json.NewEncoder(w).Encode(deviceNames)

}

func (api *Api) ClearData(w http.ResponseWriter, r *http.Request) {
	// delete all data from table
	_, err := api.Db.Query("DELETE FROM devices")
	if err != nil {
		res := &HttpRes{
			StatusCode: 400,
		}

		json.NewEncoder(w).Encode(res)
		checkErr(err)
	}
	res := &HttpRes{
		StatusCode: 200,
	}
	json.NewEncoder(w).Encode(res)
}