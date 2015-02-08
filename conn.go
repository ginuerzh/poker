package poker

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	readWait = 15 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 1024
)

type Conn struct {
	ws   *websocket.Conn
	send chan interface{}
}

func NewConn(ws *websocket.Conn, sendBuffer int) *Conn {
	conn := &Conn{
		ws:   ws,
		send: make(chan interface{}, sendBuffer),
	}
	go conn.writePump()

	return conn
}

// write writes a message with the given message type and payload.
func (c *Conn) write(mt int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(mt, payload)
}

func (c *Conn) ReadJSON(v interface{}) error {
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	if err := c.ws.ReadJSON(v); err != nil {
		log.Println(err)
		return err
	}

	if b, err := json.Marshal(v); err == nil {
		fmt.Println(">>>", time.Now().Format("15:04:05"), string(b))
	}

	return nil
}

func (c *Conn) WriteJSON(v interface{}) error {
	select {
	case c.send <- v:
		return nil
	default:
		return errors.New("buffer full")
	}
}

func (c *Conn) Close() {
	close(c.send)
	c.send = nil
}

// writePump pumps messages from the hub to the websocket connection.
func (c *Conn) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.ws.WriteJSON(message); err != nil {
				log.Println(err)
				return
			}
			if b, err := json.Marshal(message); err == nil {
				fmt.Println("<<<", time.Now().Format("15:04:05"), string(b))
			}
		case <-ticker.C:
			c.ws.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}
