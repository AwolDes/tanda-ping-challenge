package pings

import (
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	"time"
	"database/sql"
	"strings"
)

func checkErr(w http.ResponseWriter, err error) {
	switch {
		case err == sql.ErrNoRows:
			res := HttpRes{Status: 200, Description: "No rows found"}

			json.NewEncoder(w).Encode(res)
			panic(err)
		case err != nil: 
			res := HttpRes{Status: 500, Description: "Error with this request"}
	
			json.NewEncoder(w).Encode(res)
			panic(err)
	}
}

func convertTime(timestamp string) int64 {
	timestampInt, _ := strconv.ParseInt(timestamp, 10, 64)
	// if timestampInt == 0 then we were given a string
	if timestampInt == 0 {
		// The format for the date
		layout := "2006-01-02"
		date := strings.Fields(timestamp)
		// gets the date format of YYYY-MM-DD from string [2016-02-22] 10:00:00 +1000 AEST
		dateFormat := date[0]
		t, _ := time.Parse(layout, dateFormat)
		// round the day to 00:00:00
		t.Truncate(24*time.Hour)
		unix := t.Unix()
		return unix
	} else {
		return timestampInt
	}
}

func (api *Api) CreateDevicePing(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	deviceName := params["deviceId"]
	timestamp := params["timestamp"]
	timestampInt, _ := strconv.ParseInt(timestamp, 10, 64)
	
	if timestampInt == 0 {
		res := &HttpRes{
			Status: 200,
		}

		json.NewEncoder(w).Encode(res)
	}

	stmt, err := api.Db.Prepare("INSERT devices SET device_name=?,timestamp=?")
	checkErr(w, err)
	
	_, err = stmt.Exec(deviceName, timestampInt)
	checkErr(w, err)

	resSuccess := &HttpRes{
		Status: 200,
	}

	json.NewEncoder(w).Encode(resSuccess)

}

func (api *Api) GetDeviceOnDate(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	deviceName := params["deviceId"]
	timestamp := params["date"]

	timestampUnix := convertTime(timestamp)

	toTime := time.Unix(timestampUnix,0)
	toTime = toTime.AddDate(0, 0, 1)
	toTimeUnix := convertTime(toTime.String())

	rows, queryErr := api.Db.Query("SELECT device_name, timestamp FROM devices WHERE device_name=? AND timestamp>=? AND timestamp<?", deviceName, timestampUnix, toTimeUnix)
	checkErr(w, queryErr)
	defer rows.Close()

	devicePings := make([]int, 0)

	for rows.Next() {
		var deviceName string
		var timestamp int
		rows.Scan(&deviceName, &timestamp)
		devicePings = append(devicePings, timestamp)
	}

	json.NewEncoder(w).Encode(devicePings)
}

func (api *Api) GetDeviceDateRange(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	deviceName := params["deviceId"]
	fromDate := params["from"]
	toDate := params["to"]

	fromDateUnix := convertTime(fromDate)
	toDateUnix := convertTime(toDate)

	rows, queryErr := api.Db.Query("SELECT device_name, timestamp FROM devices WHERE device_name=? AND timestamp>=? AND timestamp<?", deviceName, fromDateUnix, toDateUnix)
	checkErr(w, queryErr)
	defer rows.Close()

	devicePings := make([]int, 0)

	for rows.Next() {
		var deviceName string
		var timestamp int
		rows.Scan(&deviceName, &timestamp)
		devicePings = append(devicePings, timestamp)
	}

	json.NewEncoder(w).Encode(devicePings)

}

func (api *Api) GetAllDevicesOnDate(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	timestamp := params["date"]

	timestampUnix := convertTime(timestamp)

	toTime := time.Unix(timestampUnix,0)
	toTime = toTime.AddDate(0, 0, 1)
	toTimeUnix := convertTime(toTime.String())

	rows, queryErr := api.Db.Query("SELECT device_name, timestamp FROM devices WHERE timestamp>=? AND timestamp<?", timestampUnix, toTimeUnix)
	checkErr(w, queryErr)
	defer rows.Close()

	deviceMap := make(map[string][]int)

	for rows.Next() {
		var deviceName string
		var timestamp int
		rows.Scan(&deviceName, &timestamp)
		deviceMap[deviceName] = append(deviceMap[deviceName], timestamp)
	}

	json.NewEncoder(w).Encode(deviceMap)
}
func (api *Api) GetAllDevicesInDateRange(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	fromDate := params["from"]
	toDate := params["to"]

	fromDateUnix := convertTime(fromDate)
	toDateUnix := convertTime(toDate)

	rows, queryErr := api.Db.Query("SELECT device_name, timestamp FROM devices WHERE timestamp>=? AND timestamp<?", fromDateUnix, toDateUnix)
	checkErr(w, queryErr)
	defer rows.Close()

	deviceMap := make(map[string][]int)
	
	for rows.Next() {
		var deviceName string
		var timestamp int
		rows.Scan(&deviceName, &timestamp)
		deviceMap[deviceName] = append(deviceMap[deviceName], timestamp)
	}

	json.NewEncoder(w).Encode(deviceMap)
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

	json.NewEncoder(w).Encode(deviceNames)

}

func (api *Api) ClearData(w http.ResponseWriter, r *http.Request) {
	// delete all data from table
	_, err := api.Db.Query("DELETE FROM devices")
	checkErr(w, err)
	res := &HttpRes{
		Status: 200,
	}
	json.NewEncoder(w).Encode(res)
}