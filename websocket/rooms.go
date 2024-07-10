package websocket

import "Umeet/utils"

var Rooms = room{utils.NewSafeMap(8)}

//map[string]*ClientsManager{}
type room struct {
	rooms *utils.SafeMap
}

func (r *room) Get(key any) *ClientsManager {
	t, ok := r.rooms.Get(key)
	if !ok {
		return nil
	}
	cm, _ := t.(*ClientsManager)
	return cm
}

func (r *room) Set(key, value any) {
	r.rooms.Set(key, value)
}
func (r *room) Delete(key any) {
	r.rooms.Delete(key)
}
func (r *room) Len() int {
	return r.rooms.Len()
}

// var room = map[string]*ClientsManager{}
