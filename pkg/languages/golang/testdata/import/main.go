package main

import (
	"example.com/a/v5"
	"example.com/b"
	"example.com/c-go.v5"
	"example.com/go-d"

	e "example.com/foo"
)

func m() {
	a.Test()
	b.Test()
	c.Test()
	d.Test()
	e.Test()
	other.Test()
}
