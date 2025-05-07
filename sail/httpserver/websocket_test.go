package httpserver

import (
	"github.com/gorilla/websocket"
	"testing"
)

func TestWrapWebsocketHandler(t *testing.T) {
	t.Run("WrapWebsocketHandler-NonValue", func(t *testing.T) {
		WrapWebsocketHandler(nil, nil)
	})

	t.Run("WrapWebsocketHandler", func(t *testing.T) {
		ws := &websocket.Conn{}
		handler := func(ws *websocket.Conn) {}
		WrapWebsocketHandler(ws, handler)
	})
}
