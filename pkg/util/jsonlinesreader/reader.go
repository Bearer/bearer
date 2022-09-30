package jsonlinesreader

// Reader is bufio.Scanner implementation for jsonlines
type Reader struct {
	path string
	data interface{}
	text string
}

// New creates a new report Reader
func New(path string) (*Reader, error) {
	// opens file for reading
	// creates buffio scanner
	return &Reader{}, nil
}

// Next scans next jsonline returning true upon reaching end
func (reader *Reader) Next() (end bool) {
	// reads next from bufio.scanner and saves its data returns true upon end
	return false
}

// Data returns most recent decoded interface
func (reader *Reader) Data() interface{} {
	return ""
}

// StringData returns most recent read jsonline as json string
func (reader *Reader) Text() string {
	return ""
}

// Close performs cleanup by closing file handler
func (reader *Reader) Close() {
	// close underlying file handle
}
