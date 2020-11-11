package websocket

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type StatisticFn = func() ([]byte, error)

type Client struct {
	upgrader websocket.Upgrader
	statFn   StatisticFn
	conn     *websocket.Conn
}

func New(statFn StatisticFn) *Client {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	return &Client{
		upgrader: upgrader,
		statFn:   statFn,
	}
}

func (c *Client) Upgrade(w http.ResponseWriter, r *http.Request) error {
	c.upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	conn, err := c.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}

	c.conn = conn

	return nil
}

func (c *Client) Write() {
	if c.conn == nil {
		return
	}

	ticker := time.NewTicker(5 * time.Second)
	for t := range ticker.C {
		fmt.Printf("updating stats : %+v\n", t)

		data, err := c.statFn()
		if err != nil {
			continue
		}

		if err := c.conn.WriteMessage(websocket.TextMessage, data); err != nil {
			continue
		}
	}
}
