package websocket

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// 限制为512byte时，汉字可以输入170个
// 汉字及其符号:3byte
// 空格, 英文及其符号:1byte
// emoji:4byte
const (
	writeWait = 10 * time.Second
	readWait  = 60 * time.Second
	// maxMessageSize = 512
	pingPeriod = (readWait * 9) / 10
)

type Client struct {
	uid  string
	name string
	cm   *ClientsManager
	conn *websocket.Conn
	send chan *Message
}

type Message struct {
	Uid         string `json:"uid"`
	NickName    string `json:"nickname"`
	Content     string `json:"content"`
	ImgName     string `json:"name"`
	MessageType int    `json:"messageType"`
	Size        int64  `json:"size"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  2048, //2KB
	WriteBufferSize: 2048,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// 从客户端读取数据
func (c *Client) Reader() {
	defer func() {
		c.cm.unregister <- c
		c.conn.Close()
	}()
	// c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(readWait))
	//pong处理函数就是在收到pong后可以做出反应，比如可以在收到pong后刷新下次读取的最后期限
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(readWait)); return nil })
	for {
		t, message, err := c.conn.ReadMessage()
		if err != nil {
			// if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			// 	log.Printf("error: %v", err)
			// }
			break
		}
		c.cm.broadcast <- &Message{
			Uid:         c.uid,
			NickName:    c.name,
			Content:     string(message),
			MessageType: t,
		}
	}
}

// 将数据发送
func (c *Client) Writer() {
	//心跳检测从被创建起固定每一段时间发送一次ping,
	heartbeats := time.NewTicker(pingPeriod)
	defer func() {
		heartbeats.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.conn.WriteJSON(message); err != nil {
				return
			}
		case <-heartbeats.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// 用户连接
func ServeWs(cm *ClientsManager, w http.ResponseWriter, r *http.Request, name, rid, uid string) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{cm: cm, conn: conn, send: make(chan *Message, 16), name: name, uid: uid}
	client.cm.register <- client

	go client.Writer()
	go client.Reader()
	// time.Sleep(time.Second * 10)
	// fmt.Println("over")
	// delete(Rooms, "1")
	// close(cm.Broadcast)

}
