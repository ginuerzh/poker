package poker

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"log"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	readWait = 60 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 1024
)

type Conn struct {
	ws   *websocket.Conn
	send chan []byte
}

func NewConn(ws *websocket.Conn, sendBuffer int) *Conn {
	conn := &Conn{
		ws:   ws,
		send: make(chan []byte, sendBuffer),
	}
	go conn.writePump()

	return conn
}

// write writes a message with the given message type and payload.
func (c *Conn) write(mt int, payload []byte) error {
	return c.ws.WriteMessage(mt, payload)
}

func (c *Conn) ReadJSON(v interface{}) error {
	return c.ReadJSONTimeout(v, 0)
}

func (c *Conn) ReadJSONTimeout(v interface{}, timeout time.Duration) error {
	if timeout > 0 {
		c.ws.SetReadDeadline(time.Now().Add(timeout))
	} else {
		c.ws.SetReadDeadline(time.Time{})
	}

	return c.readJson(v)
}

func (c *Conn) readJson(v interface{}) error {
	_, r, err := c.ws.NextReader()
	if err != nil {
		return err
	}
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	fmt.Println(">>>", time.Now().Format("15:04:05"), string(b))

	return json.Unmarshal(b, v)
}

func (c *Conn) WriteJSON(v interface{}) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	select {
	case c.send <- b:
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
			c.ws.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.write(websocket.TextMessage, message); err != nil {
				log.Println(err)
				return
			}

			fmt.Println("<<<", time.Now().Format("15:04:05"), string(message))
		case <-ticker.C:
			c.ws.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}
