package pings

import (
	"time"
	"net/http"
	"database/sql"
	"strings"
	"encoding/json"
	"strconv"
)


func writeJson(w http.ResponseWriter, res interface{}) {
	w.WriteHeader(http.StatusOK)
	
	json.NewEncoder(w).Encode(res)
}

func checkErr(w http.ResponseWriter, err error) {
	switch {
		case err == sql.ErrNoRows:
			res := HttpRes{Status: 200, Description: "No rows found"}
			writeJson(w, res)
			panic(err)
		case err != nil: 
			res := HttpRes{Status: 500, Description: "Error with this request"}
			w.WriteHeader(500)
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

func selectDevice(api *Api, w http.ResponseWriter, deviceName string, fromDate int64, toDate int64) []int {
	rows, queryErr := api.Db.Query("SELECT device_name, timestamp FROM devices WHERE device_name=? AND timestamp>=? AND timestamp<?", deviceName, fromDate, toDate)
	checkErr(w, queryErr)
	defer rows.Close()

	devicePings := make([]int, 0)

	for rows.Next() {
		var deviceName string
		var timestamp int
		rows.Scan(&deviceName, &timestamp)
		devicePings = append(devicePings, timestamp)
	}
	return devicePings
}

func addDay(currentUnixTime int64) int64 {
	toDate := time.Unix(currentUnixTime,0)
	toDate = toDate.AddDate(0, 0, 1)
	toDateUnix := convertTime(toDate.String())
	return toDateUnix
}
