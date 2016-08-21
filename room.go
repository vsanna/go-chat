package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type room struct {
	// 他のクライアントに転送するためのメッセージを保持する
	forward chan []byte
	// チャットルームに参加しようとしているクライアントを保持する
	join chan *client
	// チャットルームから退室しようとするクライアントを保持する
	leave chan *client
	// 在室するクライアントを保持する
	clients map[*client]bool
}

func newRoom() *room {
	return &room{
		forward: make(chan []byte),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
	}
}

func (r *room) run() {
	for {
		select {
		// r.joinから取り出せるのであれば追加の処理
		case client := <-r.join:
			r.clients[client] = true
		// r.leaveから取り出せるのであれば削除の処理
		case client := <-r.leave:
			delete(r.clients, client)
			close(client.send)
		// 転送するべきメッセージがあれば転送の処理
		// クライアントのwriteメソッドが拾ってブラウザに転送する
		case msg := <-r.forward:
			for client := range r.clients {
				select {
				case client.send <- msg:
					// メッセージを昇進
				default:
					// 送信に失敗
					delete(r.clients, client)
					close(client.send)
				}
			}
		}
	}
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  socketBufferSize,
	WriteBufferSize: socketBufferSize,
}

// チャットルームをHTTPハンドラにする
func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// websocketを利用するには、websocket.Upgrader.upgradeを利用してコネクションを首都kスうる
	// 成功したらclientを生成数rのでこれをroomに渡す
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}
	client := &client{
		socket: socket,
		send:   make(chan []byte, messageBufferSize),
		room:   r,
	}
	r.join <- client
	//クライアントの退室時の処理
	defer func() { r.leave <- client }()
	go client.write()
	// メインで実行するため接続は保持され続ける. client.read()は無限ループ
	client.read()
}
