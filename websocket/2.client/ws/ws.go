package ws

import (
	"crypto/tls"
	"errors"
	"github.com/8zhiniao/public/log"
	"github.com/gorilla/websocket"
	"github.com/jpillora/backoff"
	"net/http"
	"net/url"
	"sync"
	"time"
)

var ErrNotConnected = errors.New("websocket: not connected")

type WebSocketClient struct {
	Url       string
	ReqHeader http.Header

	Proxy            func(*http.Request) (*url.URL, error)
	TLSClientConfig  *tls.Config   //默认为空
	HandshakeTimeout time.Duration //默认为2秒
	Dialer           *websocket.Dialer

	HttpResponse *http.Response
	DialErr      error
	Conn         *websocket.Conn

	//是否显示重连日志消息
	NonVerbose       bool
	ConnectionStatus bool
	Mutex            sync.RWMutex

	// 设置心跳时间间隔
	KeepAliveIntervalTime time.Duration

	ReconnectIntervalMin    time.Duration
	ReconnectIntervalMax    time.Duration
	ReconnectIntervalFactor float64
	ReconnectIntervalJitter bool
}

type KeepAliveStatus struct {
	LastActiveTime time.Time
	Mutex          sync.RWMutex
}

func NewDefaultWebsocketClient() *WebSocketClient {

	wsClient := &WebSocketClient{}

	wsClient.SetDefaultProxy()
	wsClient.SetTLSClientConfig(&tls.Config{})
	wsClient.SetDefaultHandshakeTimeout()
	wsClient.SetDefaultDialer(wsClient.GetTLSClientConfig(), wsClient.GetHandshakeTimeout())

	wsClient.SetDefaultKeepAliveIntervalTime()

	wsClient.SetDefaultReconnectIntervalMin()
	wsClient.SetDefaultReconnectIntervalMax()
	wsClient.SetDefaultReconnectIntervalFactor()
	wsClient.SetDefaultReconnectIntervalJitter()

	return wsClient
}

func NewWebsocketClient() {

}

// Dial creates a new client connection.
// The URL url specifies the host and request URI. Use requestHeader to specify
// the origin (Origin), subprotocols (Sec-WebSocket-Protocol) and cookies
// (Cookie). Use GetHTTPResponse() method for the response.Header to get
// the selected subprotocol (Sec-WebSocket-Protocol) and cookies (Set-Cookie).
func (ws *WebSocketClient) Dial(Url string, reqHeader http.Header) {
	url, err := ws.VerifyUrl(Url)

	if err != nil {
		log.Fatalf("Dial: %v", err)
	}

	ws.SetUrl(url)
	ws.SetReqHeader(reqHeader)

	// 连接服务端
	go ws.Connect()

	time.Sleep(ws.GetHandshakeTimeout())

}

func (ws *WebSocketClient) Connect() {

	IntervalController := ws.GetBackOffController()

	for {

		NextIntervalTime := IntervalController.Duration()

		conn, response, err := websocket.DefaultDialer.Dial(ws.Url, ws.ReqHeader)

		// 则对客户端连接属性进行配置
		ws.Mutex.Lock()
		ws.HttpResponse = response
		ws.DialErr = err
		ws.Conn = conn
		if err == nil {
			ws.ConnectionStatus = true
		} else {
			ws.ConnectionStatus = false
		}
		ws.Mutex.Unlock()

		// 如果连接成功，则打印连接成功信息，并配置心跳及断线重连机制。

		if err == nil {
			log.Info("success connect to websocket server !")

			// 如果keepalive设置的心跳检时间间隔为0，则只初始化连接一次，儿不执行心跳检测及断线重连机制，如果设置的时间不为0，
			//则进行心跳检测并进行断线重连

			if ws.GetKeepAliveIntervalTime() != 0 {
				ws.Keepalive()
			}

			return

		}

		// 如果连接不成功,则睡眠一定的时间，然后继续重新连接，如果长时间连接不成功，则增加连接时间间隔到2分钟。
		log.Error("connect to websocket server not success，will sleep ", NextIntervalTime, " and try connect again !")
		time.Sleep(NextIntervalTime)

		continue

	}

}

func (ws *WebSocketClient) Keepalive() {

	ka := &KeepAliveStatus{}
	ticker := time.NewTicker(ws.GetKeepAliveIntervalTime() * time.Second)

	ws.Mutex.Lock()
	ws.Conn.SetPongHandler(func(msg string) error {
		ka.SetLastActiveTime()
		return nil
	})
	ws.Mutex.Unlock()

	go func() {
		defer ticker.Stop()
		for {

			if ws.GetConnectionStatus() == true {
				continue
			}

			if err := ws.WriteControlPingMessage(); err != nil {

				log.Error("keepalive ping remote server failed !")

			}

			<-ticker.C
			if time.Since(ka.GetLastActiveTime()) > ws.GetKeepAliveIntervalTime() {
				ws.CloseAndReconnect()
				return
			}
		}

	}()

}

// 构造函数==============================================

func (ws *WebSocketClient) SetUrl(url string) {
	ws.Mutex.Lock()
	defer ws.Mutex.Unlock()
	ws.Url = url
}

func (ws *WebSocketClient) GetUrl() string {
	ws.Mutex.Lock()
	defer ws.Mutex.Unlock()
	return ws.Url
}

func (ws *WebSocketClient) SetReqHeader(ReqHeader http.Header) {
	ws.Mutex.Lock()
	defer ws.Mutex.Unlock()

	ws.ReqHeader = ReqHeader
}

func (ws *WebSocketClient) GetReqHeader() http.Header {
	ws.Mutex.Lock()
	defer ws.Mutex.Unlock()

	return ws.ReqHeader
}

func (ws *WebSocketClient) SetDefaultHandshakeTimeout() {
	ws.Mutex.Lock()
	defer ws.Mutex.Unlock()

	if ws.HandshakeTimeout == 0 {
		ws.HandshakeTimeout = 2 * time.Second
	}
}

func (ws *WebSocketClient) GetHandshakeTimeout() time.Duration {
	ws.Mutex.Lock()
	defer ws.Mutex.Unlock()

	return ws.HandshakeTimeout
}

func (ws *WebSocketClient) SetDefaultKeepAliveIntervalTime() {
	ws.Mutex.RLock()
	defer ws.Mutex.RUnlock()

	ws.KeepAliveIntervalTime = 10
}

func (ws *WebSocketClient) GetKeepAliveIntervalTime() time.Duration {
	ws.Mutex.RLock()
	defer ws.Mutex.RUnlock()

	return ws.KeepAliveIntervalTime
}

func (ws *WebSocketClient) SetDefaultReconnectIntervalMin() {
	ws.Mutex.Lock()
	defer ws.Mutex.Unlock()

	if ws.ReconnectIntervalMin == 0 {
		ws.ReconnectIntervalMin = 2 * time.Second
	}
}

func (ws *WebSocketClient) GetReconnectIntervalMin() time.Duration {
	ws.Mutex.Lock()
	defer ws.Mutex.Unlock()
	return ws.ReconnectIntervalMin
}

func (ws *WebSocketClient) SetDefaultReconnectIntervalMax() {
	ws.Mutex.Lock()
	defer ws.Mutex.Unlock()

	if ws.ReconnectIntervalMax == 0 {
		ws.ReconnectIntervalMax = 120 * time.Second
	}
}

func (ws *WebSocketClient) GetReconnectIntervalMax() time.Duration {
	ws.Mutex.Lock()
	defer ws.Mutex.Unlock()
	return ws.ReconnectIntervalMax
}

func (ws *WebSocketClient) SetDefaultReconnectIntervalFactor() {
	ws.Mutex.Lock()
	defer ws.Mutex.Unlock()

	if ws.ReconnectIntervalFactor == 0 {
		ws.ReconnectIntervalFactor = 2
	}
}

func (ws *WebSocketClient) GetReconnectIntervalFactor() float64 {
	ws.Mutex.Lock()
	defer ws.Mutex.Unlock()
	return ws.ReconnectIntervalFactor
}

func (ws *WebSocketClient) SetDefaultReconnectIntervalJitter() {
	ws.Mutex.Lock()
	defer ws.Mutex.Unlock()

	if ws.ReconnectIntervalJitter == true {
		ws.ReconnectIntervalJitter = false
	}
}

func (ws *WebSocketClient) GetReconnectIntervalJitter() bool {
	ws.Mutex.Lock()
	defer ws.Mutex.Unlock()
	return ws.ReconnectIntervalJitter
}

func (ws *WebSocketClient) GetBackOffController() *backoff.Backoff {
	ws.Mutex.RLock()
	defer ws.Mutex.RUnlock()

	return &backoff.Backoff{
		Min:    ws.ReconnectIntervalMin,
		Max:    ws.ReconnectIntervalMax,
		Factor: ws.ReconnectIntervalFactor,
		Jitter: false,
	}
}

func (ws *WebSocketClient) SetDefaultProxy() {
	ws.Mutex.Lock()
	defer ws.Mutex.Unlock()

	if ws.Proxy == nil {
		ws.Proxy = http.ProxyFromEnvironment
	}
}

func (ws *WebSocketClient) GetProxy() func(*http.Request) (*url.URL, error) {
	ws.Mutex.Lock()
	defer ws.Mutex.Unlock()
	return ws.Proxy
}

func (ws *WebSocketClient) SetTLSClientConfig(tlsClientConfig *tls.Config) {
	ws.Mutex.Lock()
	defer ws.Mutex.Unlock()

	ws.TLSClientConfig = tlsClientConfig
}

func (ws *WebSocketClient) GetTLSClientConfig() *tls.Config {
	ws.Mutex.Lock()
	defer ws.Mutex.Unlock()

	return ws.TLSClientConfig
}

//-------------------------------------------------

// VerifyUrl parses current url

func (ws *WebSocketClient) VerifyUrl(Url string) (string, error) {
	if Url == "" {
		return "", errors.New("dial: url cannot be empty")
	}

	url, err := url.Parse(Url)

	if err != nil {
		return "", errors.New("url: " + err.Error())
	}

	if url.Scheme != "ws" && url.Scheme != "wss" {
		return "", errors.New("url: websocket uris must start with ws or wss scheme")
	}

	if url.User != nil {
		return "", errors.New("url: user name and password are not allowed in websocket url")
	}

	return Url, nil
}

func (ws *WebSocketClient) SetDefaultDialer(tlsClientConfig *tls.Config, handshakeTimeout time.Duration) {
	ws.Mutex.Lock()
	defer ws.Mutex.Unlock()

	ws.Dialer = &websocket.Dialer{
		Proxy:            ws.Proxy,
		TLSClientConfig:  tlsClientConfig,
		HandshakeTimeout: handshakeTimeout,
	}
}

// SetConnectionStatus set the WebSocket connection status
func (ws *WebSocketClient) SetConnectionStatus(status bool) {
	ws.Mutex.RLock()
	defer ws.Mutex.RUnlock()

	ws.ConnectionStatus = status

}

// GetConnectionStatus returns the WebSocket connection status
func (ws *WebSocketClient) GetConnectionStatus() bool {
	ws.Mutex.RLock()
	defer ws.Mutex.RUnlock()

	return ws.ConnectionStatus
}

//------------------ close and reconnect----------------

func (ws *WebSocketClient) CloseAndReconnect() {
	ws.Close()
	go ws.Connect()
}

func (ws *WebSocketClient) Close() {

	ws.Conn.Close()

	ws.SetConnectionStatus(false)

}

// ----------------KeepAlive------------------

func (ka *KeepAliveStatus) SetLastActiveTime() {

	ka.Mutex.Lock()
	defer ka.Mutex.Unlock()
	ka.LastActiveTime = time.Now()

}

func (ka *KeepAliveStatus) GetLastActiveTime() time.Time {

	ka.Mutex.Lock()
	defer ka.Mutex.Unlock()

	return ka.LastActiveTime

}

func (ws *WebSocketClient) WriteControlPingMessage() error {

	ws.Mutex.Lock()
	defer ws.Mutex.Unlock()
	return ws.Conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(10*time.Second))

}

// ------------------------------------

/*
使用客户端读写数据
*/

func (ws *WebSocketClient) WriteMessage(messageType int, data []byte) error {
	err := ErrNotConnected
	if ws.GetConnectionStatus() {
		ws.Mutex.Lock()
		err = ws.Conn.WriteMessage(messageType, data)
		ws.Mutex.Unlock()
		if err != nil {
			ws.CloseAndReconnect()
		}
	}
	return err
}

func (ws *WebSocketClient) ReadMessage() (messageType int, message []byte, err error) {
	err = ErrNotConnected
	if ws.GetConnectionStatus() {
		messageType, message, err = ws.Conn.ReadMessage()
		if err != nil {
			ws.CloseAndReconnect()
		}
	}

	return
}

// WriteJSON writes the JSON encoding of v to the connection.
//
// See the documentation for encoding/json Marshal for details about the
// conversion of Go values to JSON.
//
// If the connection is closed ErrNotConnected is returned
func (ws *WebSocketClient) WriteJSON(v interface{}) error {
	err := ErrNotConnected
	if ws.GetConnectionStatus() {
		ws.Mutex.Lock()
		err = ws.Conn.WriteJSON(v)
		ws.Mutex.Unlock()
		if err != nil {
			ws.CloseAndReconnect()
		}
	}

	return err
}

// ReadJSON reads the next JSON-encoded message from the connection and stores
// it in the value pointed to by v.
//
// See the documentation for the encoding/json Unmarshal function for details
// about the conversion of JSON to a Go value.
//
// If the connection is closed ErrNotConnected is returned
func (ws *WebSocketClient) ReadJSON(v interface{}) error {
	err := ErrNotConnected
	if ws.GetConnectionStatus() {
		err = ws.Conn.ReadJSON(v)
		if err != nil {
			ws.CloseAndReconnect()
		}
	}

	return err
}
