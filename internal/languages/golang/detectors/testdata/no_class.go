package foo

import "fmt"

func main() {
	user := User{
		email: "foo@bar.com",
	}

	fmt.Println(user.email)
}
