package main

import "strings"

type Command struct {
	Name string
	Args []string
}

func ParseCommand(msg string) (c Command) {
	if len(msg) < 2 {
		return
	}
	if msg[0] != '/' {
		return
	}
	s := strings.Split(msg[1:], " ")
	if len(s) == 0 {
		return
	}
	c.Name = s[0]
	if len(s) > 1 {
		c.Args = s[1:]
	}
	for key := range c.Args {
		c.Args[key] = trimString(c.Args[key])
	}
	return
}
