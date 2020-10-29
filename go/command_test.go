package main

import (
	"reflect"
	"testing"
)

func TestParseCommand(t *testing.T) {
	want1 := Command{
		Name: "login",
		Args: []string{"a", "b"},
	}
	if get := ParseCommand("/login a b"); !reflect.DeepEqual(want1, get){
		t.Errorf("want %+v, get %+v", want1, get)
	}

	want2 := Command{
		Name: "help",
	}
	if get := ParseCommand("/help\n"); !reflect.DeepEqual(want2, get){
		t.Errorf("want %+v, get %+v", want2, get)
	}

	want3 := Command{
		Name: "reg",
	}
	if get := ParseCommand("/reg\n"); !reflect.DeepEqual(want3, get){
		t.Errorf("want %+v, get %+v", want3, get)
	}
}
