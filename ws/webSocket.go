package ws

import (
	"VolumePrice/helper"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"os"
	"os/signal"
	"time"
)

type SubscribeReq struct {
	Type       string   `json:"type"`
	ProductIds []string `json:"product_ids"`
	Channels   []string `json:"channels"`
}

type WebSocket struct {
	conn *websocket.Conn
}

func (w WebSocket) SubscribeAndRead(matches chan helper.Match) error {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	flag.Parse()
	log.SetFlags(0)

	log.Printf("connecting to coinbase")
	var err error

	w.conn, _, err = websocket.DefaultDialer.Dial("wss://ws-feed.exchange.coinbase.com", nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	subscribeReq := SubscribeReq{
		Type:       "subscribe",
		ProductIds: []string{"BTC-USD", "ETH-USD", "ETH-BTC"},
		Channels:   []string{"matches"},
	}
	data, _ := json.Marshal(subscribeReq)
	err = w.conn.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		fmt.Println(err)
	}

	go func() {
		defer w.conn.Close()
		defer close(matches)
		for {
			_, message, err := w.conn.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				break
			}
			var match helper.Match
			err = json.Unmarshal(message, &match)
			if err != nil {
				fmt.Println(err)
				break
			}
			matches <- match
		}
	}()
	done := make(chan struct{})
	go func() {
		for {
			select {
			case <-done:
				return
			case <-interrupt:
				log.Println("interrupt")
				// Cleanly close the connection by sending a close message and then
				// waiting (with timeout) for the server to close the connection.
				err := w.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
				if err != nil {
					log.Println("write close:", err)
					return
				}
				select {
				case <-done:
				case <-time.After(time.Second):
				}
				return
			}
		}
	}()

	return nil
}
