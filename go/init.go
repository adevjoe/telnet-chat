package main

import "log"

func init() {
	log.Printf("init data...")
	_ = register("admin", "admin", ADMIN)
	_, _ = createRoomStore("demo")
	log.Printf("init data done")
}
