package domain

import "time"

type ClientVersion string

const (
	V1 ClientVersion = "api.bullean_ai.com"
)

type FeatureType int

const (
	FEAT_OPEN             FeatureType = 1
	FEAT_HIGH             FeatureType = 2
	FEAT_LOW              FeatureType = 3
	FEAT_CLOSE            FeatureType = 4
	FEAT_CLOSE_PERCENTAGE FeatureType = 5
)

type ClientConfig struct {
	Version        ClientVersion `json:"client_version"`
	Name           string        `json:"client_name"`
	ApiKey         string        `json:"api_key"`
	ApiSecret      string        `json:"api_secret"`
	PrepareDataset bool          `json:"prepare_dataset"`
}

type Candle struct {
	OpenTime  *time.Time `json:"open_time"`
	Open      float64    `json:"open"`
	High      float64    `json:"high"`
	Low       float64    `json:"low"`
	Close     float64    `json:"close"`
	CloseTime *time.Time `json:"close_time"`
	Volume    float64    `json:"volume"`
	Trades    []*Trade   `json:"trades"`
}

type Trade struct {
}

type Data struct {
	Name     string    `json:"name"`
	Features []float64 `json:"feature"`
	Label    int       `json:"label"`
}

type PolicyConfig struct {
	FeatName    string      `json:"feat_name"`
	FeatType    FeatureType `json:"feat_type"`
	PolicyRange int         `json:"policy_range"`
}
