package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

var Url = "ws://127.0.0.1:8080/api/agent/ws/attributes"

func main() {

	conn, response, err := websocket.DefaultDialer.Dial(Url, http.Header{})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response)
	conn.SetPongHandler(func(appData string) error {
		fmt.Println(appData)
		fmt.Println("接受到返回的pong")
		return nil
	})
	for {

		//fmt.Println(time.Now())
		fmt.Println("设置deadline")
		err1 := conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
		if err1 != nil {
			fmt.Println(err1)
		}
		fmt.Println("------------------------------------")

		fmt.Println("测试ping")
		err2 := conn.WriteMessage(websocket.PingMessage, []byte("ping------------------------------"))
		if err2 != nil {
			fmt.Println(err2)
		}
		fmt.Println("------------------------------------")

		fmt.Println("正常发送")
		err3 := conn.WriteMessage(websocket.TextMessage, []byte("正常测试"))
		if err3 != nil {
			fmt.Println(err3)
		}
		fmt.Println("------------------------------------")

		fmt.Println("测试control")
		fmt.Println(time.Now())
		err4 := conn.WriteControl(websocket.PingMessage, []byte("ping======================="), time.Now().Add(10*time.Second))
		if err4 != nil {
			fmt.Println(err4)
			fmt.Println(time.Now())
		}

		fmt.Println("------------------------------------")

		fmt.Println("返回结果")
		message, p, _ := conn.ReadMessage()
		fmt.Println(message)
		fmt.Println(string(p))
		fmt.Println("------------------------------------")
		//fmt.Println(time.Now())

		time.Sleep(10 * time.Second)

	}

}
