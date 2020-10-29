package main

type Room struct {
	ID     int
	Name   string
	User   []string
	Active bool
	Limit  int
}

const DefaultRoomLimit int = 5

func (r *Room) Exist() bool {
	return r.ID > -1 && r.Name != ""
}

func (r *Room) UserNum() int {
	return len(r.User)
}
