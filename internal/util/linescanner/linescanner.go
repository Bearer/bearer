package linescanner

import (
	"bufio"
	"io"
	"strings"
)

// Scanner iterates over lines in the given input
type Scanner struct {
	input      *bufio.Reader
	length     int
	byteOffset int
	lineNumber int
	bytes      []byte
	err        error
}

// New constructs a new line Scanner
func New(input io.Reader) *Scanner {
	return &Scanner{
		input: bufio.NewReader(input),
	}
}

func NewSize(input io.Reader, size int) *Scanner {
	return &Scanner{
		input: bufio.NewReaderSize(input, size),
	}
}

// Scan attempts to read a line and returns whether it was successful
func (scanner *Scanner) Scan() bool {
	if scanner.err != nil {
		return false
	}

	bytes, err := scanner.input.ReadBytes('\n')
	scanner.bytes = bytes
	scanner.err = err
	scanner.byteOffset += scanner.length
	scanner.length = len(bytes)
	scanner.lineNumber++

	return true
}

// Err can be called once Scan returns false to see if the scan ended due to
// an error, or whether the end of the file was reached (Err returns nil)
func (scanner *Scanner) Err() error {
	if scanner.err == io.EOF {
		return nil
	}

	return scanner.err
}

// Text returns the text of the current line
func (scanner *Scanner) Text() string {
	text := string(scanner.bytes)
	return strings.TrimRight(text, "\r\n")
}

// Bytes returns the bytes of the current line
func (scanner *Scanner) Bytes() []byte {
	return scanner.bytes
}

// LineNumber returns the 1-based offset of the current line
func (scanner *Scanner) LineNumber() int {
	return scanner.lineNumber
}

// ByteOffset returns the 0-based byte offset of the start of the current line
func (scanner *Scanner) ByteOffset() int {
	return scanner.byteOffset
}
