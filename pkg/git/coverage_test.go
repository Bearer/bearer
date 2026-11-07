package git

import (
	"errors"
	"io"
	"strings"
	"testing"
)

type coverageReadFault struct{}

func (coverageReadFault) Read([]byte) (int, error) {
	return 0, errors.New("synthetic acquisition read failure")
}
func TestCoverageTreeStreamIntegrity(t *testing.T) {
	hash := strings.Repeat("a", 40)
	for _, tc := range []struct {
		name, body string
		other, ok  bool
	}{
		{"empty", "", false, true}, {"tree", "100644 blob " + hash + "\tinput.go\x00", false, true},
		{"unterminated", "100644 blob " + hash + "\tinput.go", false, false},
		{"partial metadata", "100644 blob", false, false}, {"bad metadata", "bad\tinput.go\x00", false, false},
		{"bad hash", "100644 blob nope\tinput.go\x00", false, false},
		{"empty record", "\x00", false, false}, {"other", "input.go\x00", true, true}, {"partial other", "input.go", true, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var rows []TreeFile
			err := readTreeRecords(strings.NewReader(tc.body), tc.other, &rows)
			if (err == nil) != tc.ok {
				t.Fatalf("err=%v rows=%v", err, rows)
			}
		})
	}
	var rows []TreeFile
	if err := readTreeRecords(coverageReadFault{}, false, &rows); err == nil || err == io.EOF {
		t.Fatal("read failure lost")
	}
}
