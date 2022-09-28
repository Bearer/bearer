package datatypes

type User struct {
	name string
	// supported without type linking
	person person
	*BodyStats
}

type person struct {
	city string
}

type BodyStats struct {
	height int
	width  int
}

func (user *User) Save(personParam2 *person, personParam1 person) error {
	return nil
}

// return types of methods are not supported
func (user *User) Load(packageParam2 *models.person, packageParam1 models.person) {
	return nil, models.person{}
}

func (user User) Store() {

}

// not supported as those rarely contain pii
type staticNumber int

// interfaces are not supported as they should be supported by implementations
type Producer interface {
	doSomething()
}
