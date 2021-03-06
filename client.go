package main

import (
	"github.com/gorilla/websocket"
	"log"
)

type client struct {
	// このクライアントのためのwebsocket
	socket *websocket.Conn
	// メッセージをためておく
	send chan []byte
	room *room
}

// clientがroom.forwardにsocketの保有するメッセージを貯めこむ
func (c *client) read() {
	for {
		if _, msg, err := c.socket.ReadMessage(); err == nil {
			log.Println("receive: ", string(msg), "from: ", *c)
			c.room.forward <- msg
		} else {
			break
		}
	}
	c.socket.Close()
}

// clientがsocketを介し、ブラウザに転送する
func (c *client) write() {
	for msg := range c.send {
		log.Println("send: ", string(msg), "to: ", *c)
		if err := c.socket.WriteMessage(websocket.TextMessage, msg); err != nil {
			break
		}
	}
	c.socket.Close()
}
