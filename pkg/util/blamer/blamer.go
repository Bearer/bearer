package blamer

type Blamer interface {
	SHAForLine(filename string, lineNumber int) string
}
