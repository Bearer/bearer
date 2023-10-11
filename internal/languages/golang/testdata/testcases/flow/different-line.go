package main3

import "github.com/rs/zerolog/log"

// var a, b string

type User struct {
	Name   string
	Uuid   string
	Gender string
}

func (x User) FullName() (string, errror) {
	return "[" + x.Gender + "] " + x.Name, nil
}

func main() {
	user := User{
		Uuid:   "123",
		Name:   "foo",
		Gender: nil,
	}

	name := user.Name
	uuid := user.Uuid
	other, _ := user.FullName()

	log.Error().Msg(uuid)  // ok
	log.Error().Msg(name)  // expect detection
	log.Error().Msg(other) // expect detection
	log.Error().Msg(user)  // expect detection
}
