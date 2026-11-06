package httpserver

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// WrapWebsocketHandler wraps a WebSocket handler
//
// ws WebSocket connection instance
//
// handler WebSocket message handler
func WrapWebsocketHandler(ws *websocket.Conn, handler func(ws *websocket.Conn)) func(ginContext *gin.Context) {
	return func(ginContext *gin.Context) {
		if ws == nil {
			ws = wrappedDefaultWebsocketConn(ginContext)
		}
		go func() {
			<-ginContext.Done()
			fmt.Println("[GO-SAIL] websocket lost connection")
		}()

		if handler == nil {
			defaultWebsocketHandlerFunc(ws)
		} else {
			handler(ws)
		}
	}
}

// wrappedDefaultWebsocketConn wraps the default WebSocket connection instance
func wrappedDefaultWebsocketConn(ginContext *gin.Context) *websocket.Conn {
	var upgrade = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	ws, err := upgrade.Upgrade(ginContext.Writer, ginContext.Request, nil)
	if err != nil {
		message := fmt.Errorf("[GO-SAIL] can not start up websocket: %s", err.Error())
		panic(message)
	}

	return ws
}

// defaultWebsocketHandlerFunc default WebSocket handler
var defaultWebsocketHandlerFunc = func(ws *websocket.Conn) {
	defer func() {
		_ = ws.Close()
	}()

	for {
		messageType, message, readErr := ws.ReadMessage()
		if readErr != nil {
			fmt.Printf("[GO-SAIL] read message error => messageType: %d, message: %s, error: %s\n", messageType, string(message), readErr.Error())
			break
		}
		if string(message) == "ping" {
			_ = ws.WriteMessage(messageType, []byte("[GO-SAIL] server reply: pong"))
			continue
		}

		fmt.Printf("[GO-SAIL] read message => messageType: %d, message: %s\n", messageType, string(message))

		var reply bytes.Buffer
		reply.WriteString("[GO-SAIL] server reply: ")
		reply.Write(message)
		_ = ws.WriteMessage(messageType, reply.Bytes())
	}
}
