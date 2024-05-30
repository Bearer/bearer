package main2

import "github.com/rs/zerolog/log"

type User struct {
	Name string
}

func (x User) FullName() string {
	return x.Name
}

func main() {
	user := User{
		Name:   "foo",
		Gender: "female",
	}

	log.Error().Msg(user.Name)                        // expect detection
	log.Error().Msg(user.FullName())                  // expect detection
	log.Error().Msgf("user info %s", user.FullName()) // expect detection
}
