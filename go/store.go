package main

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"sync"
)

// Global var
var (
	GlobalUsers sync.Map
	GlobalRooms = Rooms{
		RoomMap:     sync.Map{},
		RoomNameMap: sync.Map{},
	}
	GlobalCurrentRoomsID = &CurrentRoomsID{
		ID:  -1,
		Mux: &sync.Mutex{},
	}
	GlobalClient = Clients{
		Clients: make([]*Client, 0),
		Mux:     &sync.Mutex{},
	}
)

type Clients struct {
	Clients []*Client
	Mux     *sync.Mutex
}

type Rooms struct {
	RoomMap     sync.Map
	RoomNameMap sync.Map
}

type CurrentRoomsID struct {
	ID  int
	Mux *sync.Mutex
}

func (c *CurrentRoomsID) Next() int {
	if c.Mux == nil {
		return -1
	}
	c.Mux.Lock()
	defer c.Mux.Unlock()
	c.ID += 1
	return c.ID
}

func createRoomStore(name string) (Room, error) {
	room := Room{
		ID:     GlobalCurrentRoomsID.Next(),
		Name:   name,
		Active: true,
		Limit:  DefaultRoomLimit,
	}
	r := getStoreRoomByName(name)
	if r.Exist() {
		return Room{}, errors.New("room exist")
	}
	GlobalRooms.RoomMap.Store(room.ID, room)
	GlobalRooms.RoomNameMap.Store(room.Name, room.ID)
	return room, nil
}

func getStoreRoomByName(name string) Room {
	if r, ok := GlobalRooms.RoomNameMap.Load(name); ok {
		if id, ok := r.(int); ok {
			getStoreRoomByID(id)
		}
	}
	return Room{}
}

func getStoreRoomByID(id int) Room {
	if r, ok := GlobalRooms.RoomMap.Load(id); ok {
		if room, ok := r.(Room); ok {
			return room
		}
	}
	return Room{}
}

func listRooms() []string {
	var list []string
	GlobalRooms.RoomMap.Range(func(key, value interface{}) bool {
		if room, ok := value.(Room); ok {
			if room.Active {
				list = append(list, fmt.Sprintf("%d-%s(%d/%d)",
					room.ID, room.Name, getRoomUser(room.ID), room.Limit))
			}
		}
		return true
	})
	return list
}

func joinRoom(username string, roomID int) error {
	if data, ok := GlobalRooms.RoomMap.Load(roomID); ok {
		if room, ok := data.(Room); ok {
			if !room.Exist() {
				return errors.New("room none")
			}
			room.User = append(room.User, username)
			GlobalRooms.RoomMap.Store(roomID, room)
		}
	}
	return nil
}

func leaveRoom(username string, roomID int) error {
	room := getStoreRoomByID(roomID)
	for key, u := range room.User {
		if u == username {
			room.User = append(room.User[:key], room.User[key+1:]...)
			GlobalRooms.RoomMap.Store(roomID, room)
		}
	}
	GlobalClient.Mux.Lock()
	defer GlobalClient.Mux.Unlock()
	for key, c := range GlobalClient.Clients {
		if c.CurrentRoom.ID == roomID {
			GlobalClient.Clients[key].CurrentRoom = Room{ID: -1}
		}
	}
	return nil
}

func getRoomUser(id int) int {
	r := getStoreRoomByID(id)
	return len(r.User)
}

func appendClient(c *Client) {
	GlobalClient.Mux.Lock()
	defer GlobalClient.Mux.Unlock()
	GlobalClient.Clients = append(GlobalClient.Clients, c)
}

func removeClient(uuid uuid.UUID) {
	GlobalClient.Mux.Lock()
	defer GlobalClient.Mux.Unlock()
	for key, c := range GlobalClient.Clients {
		if c.UUID == uuid {
			GlobalClient.Clients = append(GlobalClient.Clients[:key], GlobalClient.Clients[key+1:]...)
		}
	}
}

func sendRoom(username string, roomID int, msg string) {
	r := getStoreRoomByID(roomID)
	for _, c := range GlobalClient.Clients {
		if c.CurrentRoom.ID == roomID && c.CurrentRoom.Name == r.Name {
			go func(c *Client, r Room) {
				c.Message <- fmt.Sprintf("[ (%d)%s ] %s: %s\n> ", r.ID, r.Name, username, msg)
			}(c, r)
		}
	}
}

func deleteRoom(roomID int) {
	room := getStoreRoomByID(roomID)
	for _, u := range room.User {
		_ = leaveRoom(u, roomID)
	}
	GlobalRooms.RoomMap.Delete(roomID)
	GlobalRooms.RoomNameMap.Delete(room.Name)
}

func checkLogin(username string) bool {
	for _, c := range GlobalClient.Clients {
		if c.User.Username == username {
			return true
		}
	}
	return false
}
