package data

import (
	"encoding/json"
	"github.com/bullean-ai/bullean-sdk/data/domain"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"net/url"
)

type Client struct {
	conn *websocket.Conn `json:"conn"`
	Name string          `json:"name"`
}

func NewClient(config domain.ClientConfig) domain.IClient {

	u := url.URL{Scheme: "ws", Host: string(config.Version), Path: "/ws"}
	log.Printf("connecting to %s", u.String())

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), http.Header{
		"api-key":    []string{config.ApiKey},
		"api-secret": []string{config.ApiSecret},
	})
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer conn.Close()

	return &Client{conn: conn, Name: config.Name}
}

func (c Client) OnCandle(fn func(domain.Candle)) {
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			var candle domain.Candle
			_, message, err := c.conn.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				continue
			}
			err = json.Unmarshal(message, &candle)
			if err != nil {
				log.Println("marshal:", err)
				continue
			}
			fn(candle)
		}
	}()
}

func (c Client) OnReady() []domain.Candle {
	//TODO implement me
	panic("implement me")
}
