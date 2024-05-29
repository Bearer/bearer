package foo

import "fmt"

type User struct {
	firstName string
	lastName  string
	email     string
	uuid      string
}

func (x User) Name() {
	fmt.Println(x.firstName)
}
