package main

import "example.com/bar/v5"
import "example.com/foo"

import (
	baz "example.com/foo"
)

func m() {
	foo.Test()
	bar.Test()
	baz.Test()
	other.Test()
}
