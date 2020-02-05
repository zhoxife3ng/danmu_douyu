package client

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"sync"
	"time"
)

type websocketClient struct {
	scheme      string
	siteName    string
	port        int
	conn        *websocket.Conn
	close       chan struct{}
	closeOnce   sync.Once
	onCloseFunc func()
}

func NewWebsocketClient(scheme, siteName string, port int) *websocketClient {
	client := &websocketClient{
		scheme:   scheme,
		siteName: siteName,
		port:     port,
		close:    make(chan struct{}),
	}
	return client
}

func (client *websocketClient) Connect() error {

	u := url.URL{Scheme: client.scheme, Host: fmt.Sprintf("%s:%d", client.siteName, client.port)}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
		return err
	}

	client.conn = c
	return nil
}

func (client *websocketClient) SetTickerFunc(f func() bool, d time.Duration) {
	ticker := time.NewTicker(d)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-client.close:
				return
			case <-ticker.C:
				if !f() {
					return
				}
			}
		}
	}()
}

func (client *websocketClient) SendMsg(msg []byte) {
	err := client.conn.WriteMessage(websocket.TextMessage, msg)
	if err != nil {
		client.Close()
	}
}

func (client *websocketClient) OnReceive(f func([]byte) bool) {
	go func() {
		defer client.Close()
		for {
			_, msg, err := client.conn.ReadMessage()
			if err != nil {
				return
			}
			if !f(msg) {
				return
			}
		}
	}()
}

func (client *websocketClient) Close() {
	client.closeOnce.Do(func() {
		close(client.close)
		if client.conn != nil {
			_ = client.conn.Close()
		}
		if client.onCloseFunc != nil {
			client.onCloseFunc()
		}
	})
}

func (client *websocketClient) OnClose(f func()) {
	client.onCloseFunc = f
}
