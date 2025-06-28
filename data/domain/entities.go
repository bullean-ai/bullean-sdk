package domain

import "time"

const HISTORY_LIMIT = 40000

type ClientVersion string

const (
	V1 ClientVersion = "152.89.38.10:5067"
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
	Version      ClientVersion `json:"client_version"`
	Name         string        `json:"client_name"`
	ApiKey       string        `json:"api_key"`
	ApiSecret    string        `json:"api_secret"`
	StreamReqMsg StreamReqMsg  `json:"subscription"`
}

type StreamReqMsg struct {
	TypeOf        string         `json:"type_of"`
	History       bool           `json:"history"`
	HistorySize   int            `json:"history_size"`
	Subscriptions []Subscription `json:"subscriptions"`
}

type Subscription struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type ResponseType string

const (
	HISTORY    ResponseType = "history"
	NEW_CANDLE ResponseType = "new_candle"
)

type StreamResMsg struct {
	TypeOf  ResponseType `json:"type_of"`
	Candles []Candle     `json:"candle"`
	IsDone  bool         `json:"is_done"`
}
type Candle struct {
	Symbol    string     `json:"s"`
	OpenTime  *time.Time `json:"t"`
	Open      float64    `json:"o"`
	High      float64    `json:"h"`
	Low       float64    `json:"l"`
	Close     float64    `json:"c"`
	CloseTime *time.Time `json:"T"`
	Volume    float64    `json:"v"`
	Trades    []*Trade   `json:"tr"`
}

type Trade struct {
	EventType    string     `json:"e"`
	EventTime    *time.Time `json:"E"`
	Symbol       string     `json:"s"`
	TradeId      string     `json:"t"`
	Price        float64    `json:"p"`
	Quantity     float64    `json:"q"`
	TradeTime    *time.Time `json:"T"`
	IsBuyerMaker bool       `json:"m"`
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
