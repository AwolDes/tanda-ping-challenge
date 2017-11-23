package pings

import (
	"net/http"
	"github.com/gorilla/mux"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	"encoding/json"
)

func (api *Api) CreateDevicePing(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	deviceName := params["deviceId"]
	timestamp := params["timestamp"]
	timestampInt, _ := strconv.ParseInt(timestamp, 10, 64)
	
	if timestampInt == 0 {
		res := &HttpRes{
			Status: 400,
			Description: "Tried to create a ping with an ISO string, must be a unix timestamp",
		}

		w.WriteHeader(400)
		json.NewEncoder(w).Encode(res)

	} else {
		stmt, err := api.Db.Prepare("INSERT devices SET device_name=?,timestamp=?")
		checkErr(w, err)
		
		_, err = stmt.Exec(deviceName, timestampInt)
		checkErr(w, err)
	
		resSuccess := &HttpRes{
			Status: 200,
			Description: "Added device ping",
		}
	
		writeJson(w, resSuccess)
	}


}

func (api *Api) GetDeviceOnDate(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	deviceName := params["deviceId"]
	timestamp := params["date"]

	fromDateUnix := convertTime(timestamp)
	toDateUnix := addDay(fromDateUnix)

	devicePings := selectDevice(api, w, deviceName, fromDateUnix, toDateUnix)

	writeJson(w, devicePings)
}

func (api *Api) GetDeviceDateRange(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	deviceName := params["deviceId"]
	fromDate := params["from"]
	toDate := params["to"]

	fromDateUnix := convertTime(fromDate)
	toDateUnix := convertTime(toDate)

	devicePings := selectDevice(api, w, deviceName, fromDateUnix, toDateUnix)
	writeJson(w, devicePings)

}

func selectAllDevices(api *Api, w http.ResponseWriter, fromDate int64, toDate int64) map[string][]int {

	rows, queryErr := api.Db.Query("SELECT device_name, timestamp FROM devices WHERE timestamp>=? AND timestamp<?", fromDate, toDate)
	checkErr(w, queryErr)
	defer rows.Close()

	deviceMap := make(map[string][]int)

	for rows.Next() {
		var deviceName string
		var timestamp int
		rows.Scan(&deviceName, &timestamp)
		deviceMap[deviceName] = append(deviceMap[deviceName], timestamp)
	}

	return deviceMap
}

func (api *Api) GetAllDevicesOnDate(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	timestamp := params["date"]

	fromDateUnix := convertTime(timestamp)
	toDateUnix := addDay(fromDateUnix)

	deviceMap := selectAllDevices(api, w, fromDateUnix, toDateUnix)

	writeJson(w, deviceMap)
}
func (api *Api) GetAllDevicesInDateRange(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	fromDate := params["from"]
	toDate := params["to"]

	fromDateUnix := convertTime(fromDate)
	toDateUnix := convertTime(toDate)

	deviceMap := selectAllDevices(api, w, fromDateUnix, toDateUnix)
	writeJson(w, deviceMap)
}

func (api *Api) GetAllDevices(w http.ResponseWriter, r *http.Request) {

	rows, err := api.Db.Query("SELECT DISTINCT device_name FROM devices")

	checkErr(w, err)

	defer rows.Close()

	deviceNames := make([]string, 0)

	for rows.Next() {
		var deviceName string
		rows.Scan(&deviceName)
		deviceNames = append(deviceNames, deviceName)
	}

	writeJson(w, deviceNames)

}

func (api *Api) ClearData(w http.ResponseWriter, r *http.Request) {
	// delete all data from table
	_, err := api.Db.Query("DELETE FROM devices")
	checkErr(w, err)
	res := &HttpRes{
		Status: 200,
		Description: "Successfully cleared data",
	}

	writeJson(w, res)
}