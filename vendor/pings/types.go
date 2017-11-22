package pings

type Device struct {
	Id string `json:"id"`
	Timestamps []int `json:"timestamps"`
}

type HttpRes struct {
	StatusCode int `json:"statusCode"`
}

