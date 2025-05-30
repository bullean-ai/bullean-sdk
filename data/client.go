package data

import (
	"encoding/json"
	"github.com/bullean-ai/bullean-sdk/data/domain"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"net/url"
)

type client struct {
	conn    *websocket.Conn `json:"conn"`
	Name    string          `json:"name"`
	DataSet domain.IDataSet `json:"data_set"`
}

func NewClient(config domain.ClientConfig) domain.IClient {
	var msg []byte

	if config.StreamReqMsg.HistorySize > domain.HISTORY_LIMIT {
		config.StreamReqMsg.HistorySize = domain.HISTORY_LIMIT
	}

	u := url.URL{Scheme: "ws", Host: string(config.Version), Path: "/ws"}
	log.Printf("connecting to %s", u.String())

	headers := http.Header{}
	headers.Add("secret-key", config.ApiKey)
	headers.Add("secret-token", config.ApiSecret)

	conn, res, err := websocket.DefaultDialer.Dial(u.String(), headers)
	if err != nil {
		if res != nil && res.StatusCode == 401 {
			log.Println("Unauthorized: Please check your API key and secret.")
			return nil
		} else {
			log.Fatal("dial:", err)
		}
	}

	msg, err = json.Marshal(config.StreamReqMsg)
	if err != nil {
		log.Println("Error JSON Marshal subscription config:", err)
		return nil
	}
	conn.WriteMessage(websocket.BinaryMessage, msg)

	return &client{conn: conn, Name: config.Name}
}

func (c client) OnReady(fn func([]domain.Candle)) {
	done := make(chan struct{})
	var msg domain.StreamResMsg
	var candles []domain.Candle
	defer close(done)
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			continue
		}
		err = json.Unmarshal(message, &msg)
		if err != nil {
			log.Println("marshal:", err)
			continue
		}
		if msg.TypeOf == domain.HISTORY {
			for _, candle := range msg.Candles {
				candles = append([]domain.Candle{candle}, candles...)
			}
		}
		if msg.IsDone {
			break
		}
	}
	if msg.IsDone {
		fn(candles)
		candles = nil
	}
}

func (c client) OnCandle(fn func([]domain.Candle)) {
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			var msg domain.StreamResMsg
			_, message, err := c.conn.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				continue
			}
			err = json.Unmarshal(message, &msg)
			if err != nil {
				log.Println("marshal:", err)
				continue
			}
			if msg.TypeOf == domain.NEW_CANDLE {
				fn(msg.Candles)
			}
		}
	}()
}
