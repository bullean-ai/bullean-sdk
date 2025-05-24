package domain

import "time"

type ClientVersion string

const (
	V1 ClientVersion = "https://api.bullean_ai.com/v1"
)

type ClientConfig struct {
	Version   ClientVersion `json:"client_version"`
	Name      string        `json:"client_name"`
	ApiKey    string        `json:"api_key"`
	ApiSecret string        `json:"api_secret"`
}

type Candle struct {
	OpenTime  *time.Time `json:"open_time"`
	Open      float64    `json:"open"`
	High      float64    `json:"high"`
	Low       float64    `json:"low"`
	Close     float64    `json:"close"`
	CloseTime *time.Time `json:"close_time"`
	Volume    float64    `json:"volume"`
}
