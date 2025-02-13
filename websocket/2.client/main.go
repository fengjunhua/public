package main

import (
	"client/ws"
	"fmt"
	"github.com/8zhiniao/public/log"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

func main() {

	log.InitLoggerFromParams("/tmp/ws-client.log", "info", "console", 10, 2, 1, true, true)
	defer log.Sync()
	client := ws.NewDefaultWebsocketClient()
	client.Dial("ws://127.0.0.1:8001/oslog", http.Header{})

	for {

		err := client.WriteMessage(websocket.TextMessage, []byte("{\"key\": \"openevent\", \"value\": \"test\"}"))
		if err != nil {
			fmt.Println(err)
		}

		message, bytes, err1 := client.ReadMessage()
		if err1 != nil {
			fmt.Println(err1)
		}
		fmt.Println(message)
		fmt.Println(string(bytes))

		time.Sleep(10 * time.Second)

	}

}
