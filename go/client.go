package main

import (
	"fmt"
	"github.com/google/uuid"
	"log"
	"net"
)

type Client struct {
	UUID        uuid.UUID
	Conn        net.Conn
	User        User
	CurrentRoom Room
	Message     chan string
	Done        chan bool
}

func NewClient(conn net.Conn) *Client {
	c := &Client{
		UUID:        uuid.New(),
		Conn:        conn,
		User:        User{},
		CurrentRoom: Room{ID: -1},
		Done:        make(chan bool),
		Message:     make(chan string),
	}
	appendClient(c)
	return c
}

func (c *Client) Receive() (string, error) {
	msg := make([]byte, 100)
	_, err := c.Conn.Read(msg)
	if err != nil {
		return "", err
	}
	return string(msg), nil
}

func (c *Client) Send(msg string) {
	roomName := ""
	if c.CurrentRoom.Name != "" {
		roomName = fmt.Sprintf("[ (%d)%s ] ", c.CurrentRoom.ID, c.CurrentRoom.Name)
	}
	_, _ = c.Conn.Write([]byte(roomName + msg + "\n> "))
}

func (c *Client) Close() {
	c.LeaveRoom()
	c.Done <- true
	removeClient(c.UUID)
	_ = c.Conn.Close()
	log.Printf("client %s disconnect.\n", c.UUID)
}

func (c *Client) Auth() bool {
	return c.User.Username != ""
}

func (c *Client) Chat(msg string) {
	if !c.Auth() {
		return
	}
	if c.CurrentRoom.ID == -1 {
		c.Send("")
		return
	}
	sendRoom(c.User.Username, c.CurrentRoom.ID, msg)
}

func (c *Client) LeaveRoom() {
	if !c.Auth() {
		return
	}
	if c.CurrentRoom.ID == -1 {
		return
	}
	_ = leaveRoom(c.User.Username, c.CurrentRoom.ID)
}
