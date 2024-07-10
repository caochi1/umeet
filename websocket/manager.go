package websocket

import (
	"Umeet/utils"
)

func NewManager() *ClientsManager {
	return &ClientsManager{
		make(chan *Message, 16),
		make(chan *Client),
		make(chan *Client),
		make(map[*Client]struct{}),
		utils.NewSafeMap(8),
	}
}

type ClientsManager struct {
	broadcast  chan *Message
	register   chan *Client
	unregister chan *Client
	clients    map[*Client]struct{}
	member     *utils.SafeMap
	// member  map[string]string uid|nickname

}

func (cm *ClientsManager) Run() {
	for {
		select {
		case client := <-cm.register:
			cm.clients[client] = struct{}{}
		case client := <-cm.unregister:
			if _, ok := cm.clients[client]; ok {
				cm.remove(client)
			}
		case message, open := <-cm.broadcast:
			if !open {
				for client := range cm.clients {
					cm.remove(client)
				}
				return
			}
			for client := range cm.clients {
				select {
				case client.send <- message:
				default:
					cm.remove(client)
				}
			}
		}
	}
}

func (cm *ClientsManager) Get(key any) (string, bool) {
	t, ok := cm.member.Get(key)
	if !ok {
		return "", ok
	}
	memb, _ := t.(string)
	return memb, ok
}

func (cm *ClientsManager) Set(key, value any) {
	cm.member.Set(key, value)
}

func (cm *ClientsManager) Delete(key any) {
	cm.member.Delete(key)
}

func (cm *ClientsManager) Copy() map[string]string {
	newMap := make(map[string]string, cm.member.Len())
	cm.member.ForEach(func(k, v any) {
		key, _ := k.(string)
		value, _ := v.(string)
		newMap[key] = value
	})
	return newMap
}

func (cm *ClientsManager) Send(msg *Message) {
	cm.broadcast <- msg
}

func (cm *ClientsManager) Close() {
	close(cm.broadcast)
}

func (cm *ClientsManager) remove(client *Client) {
	close(client.send)
	delete(cm.clients, client)
}

// func (cm *ClientsManager) locker(f func(*Client)) func(*Client) {
// 	return func(client *Client) {
// 		cm.lock.Lock()
// 		f(client)
// 		cm.lock.Unlock()
// 	}
// }
