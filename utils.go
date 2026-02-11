package main

import (
	"github.com/godbus/dbus/v5"
	"github.com/godbus/dbus/v5/prop"
)

func newProp(value any, cb func(*prop.Change) *dbus.Error) *prop.Prop {
	return &prop.Prop{
		Value:    value,
		Writable: true,
		Emit:     prop.EmitTrue,
		Callback: cb,
	}
}
