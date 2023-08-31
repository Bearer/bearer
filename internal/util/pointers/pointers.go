package pointers

func String(value string) *string {
	return &value
}

func Int(value int) *int {
	return &value
}

func Bool(value bool) *bool {
	return &value
}
