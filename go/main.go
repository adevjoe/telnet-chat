package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
)

func main() {
	ln, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Printf("listen error: %v\n", err)
		return
	}
	log.Printf("start listening...\n")
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("%s connect fail, error: %v\n", conn.RemoteAddr(), err)
			continue
		}
		go handleConnection(NewClient(conn))
	}
}

func handleConnection(c *Client) {
	log.Printf("client %s connected.\n", c.UUID)
	c.Send("Please login. Use command '/help' to get help.")
	go func(c *Client) {
		for {
			select {
			case <-c.Done:
				return
			case msg := <-c.Message:
				_, _ = c.Conn.Write([]byte(msg))
			}
		}
	}(c)
	for {
		msg, err := c.Receive()
		if err != nil {
			if err == io.EOF {
				c.Close()
				return
			}
			log.Printf("client %s connect error: %v\n", c.UUID, err)
			return
		}
		if len(msg) > 0 {
			if processMsg(c, msg) == -1 { // connection closed
				return
			}
		}
	}
}

func processMsg(c *Client, msg string) int {
	msg = trimString(msg)
	command := ParseCommand(msg)
	if command.Name == "" { // chat
		c.Chat(msg)
		return 0
	}
	switch command.Name {
	case "login":
		if len(command.Args) != 2 {
			c.Send("Command error.")
			break
		}
		var (
			username = command.Args[0]
			password = command.Args[1]
		)
		if c.Auth() {
			c.Send("You're logged in.")
			break
		}
		if checkPassword(username, password) {
			if checkLogin(username) {
				c.Send("You're logged in other session.")
				break
			}
			c.User = getUserByUsername(username)
			c.Send("Login successful.")
		} else {
			c.Send("Login fail.")
		}
	case "reg", "register":
		if len(command.Args) != 2 {
			c.Send("Command error.")
			break
		}
		err := register(command.Args[0], command.Args[1], MEMBER)
		if err != nil {
			c.Send(fmt.Sprintf("Error: %s", err.Error()))
		} else {
			c.Send("Register successful.")
		}
	case "i":
		if !c.Auth() {
			c.Send("Need login.")
			break
		}
		c.Send(fmt.Sprintf("Username -> %s, Level -> %s", c.User.Username,
			c.User.Level))
	case "logout":
		if !c.Auth() {
			c.Send("Need login.")
			break
		}
		c.Close()
		c.User = User{}
		c.Send("Logout successful.")
	case "rooms":
		if !c.Auth() {
			c.Send("Need login.")
			break
		}
		c.Send(fmt.Sprintf("rooms: %v", listRooms()))
	case "join":
		if !c.Auth() {
			c.Send("Need login.")
			break
		}
		if len(command.Args) != 1 {
			c.Send("Command error.")
			break
		}
		roomID, err := strconv.Atoi(command.Args[0])
		if err != nil {
			c.Send("room error.")
			break
		}
		room := getStoreRoomByID(roomID)
		if room.Name != "" && room.Active {
			count := getRoomUser(roomID)
			if count >= room.Limit {
				c.Send("This room is full up.")
				break
			}
			err := joinRoom(c.User.Username, roomID)
			if err != nil {
				c.Send(err.Error())
				break
			}
		} else {
			c.Send("Room none.")
			break
		}
		c.CurrentRoom = room
		c.Send(fmt.Sprintf("Welcome to room %d.", roomID))
	case "leave":
		if !c.Auth() {
			c.Send("Need login.")
			break
		}
		if c.CurrentRoom.ID == -1 {
			c.Send("You are not in the room.")
			break
		}
		id := c.CurrentRoom.ID
		_ = leaveRoom(c.User.Username, c.CurrentRoom.ID)
		c.Send(fmt.Sprintf("Leave room %d.", id))
	case "create":
		if !c.Auth() || !c.User.IsAdmin() {
			c.Send("You are not allow to do that.")
			break
		}
		name := command.Args[0]
		room, err := createRoomStore(name)
		if err != nil {
			c.Send(err.Error())
			break
		}
		c.Send(fmt.Sprintf("You created a room. (ID: %d, Name: %s)",
			room.ID, room.Name))
	case "del", "delete":
		if !c.Auth() || !c.User.IsAdmin() {
			c.Send("You are not allow to do that.")
			break
		}
		roomID, err := strconv.Atoi(command.Args[0])
		if err != nil {
			c.Send("room error.")
			break
		}
		deleteRoom(roomID)
		c.Send(fmt.Sprintf("Room %d deleted.", roomID))
	case "exit":
		c.Close()
		return -1
	case "help":
		c.Send(`Chat Room.
Usage:
  /command [args...]
Commands:
   register <username> <password>    create an account
   login <username> <password>       login with your username and password
   logout                            logout your account
   i                                 get my info
   rooms                             list rooms
   join <room_id>                    join a room
   leave <room_id>                   leave a room
   exit                              exit this connection
   help                              get help tip
`)
	default:
		c.Send("Unknown command. Tap /help to get help.")
	}
	return 0
}
