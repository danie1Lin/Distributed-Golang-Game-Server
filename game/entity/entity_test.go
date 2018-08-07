package entity

import (
	"reflect"
	"testing"
)

func TestRegisterAndCall(*testing.T) {
	RegisterEnitity(&Player{})
	id := eManager.CreateEnitity("Player", true)
	eManager.Call("Player", id, "Say", reflect.ValueOf("Yo"))
	id = eManager.CreateEnitity("Player", true)
	eManager.Call("Player", id, "Say", reflect.ValueOf("HI"))
	id = eManager.CreateEnitity("Player", true)
	eManager.Call("Player", id, "Say", reflect.ValueOf("HI"))
}
