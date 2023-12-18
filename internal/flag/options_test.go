package flag

import (
	"testing"
)

func Test_getStringSlice(t *testing.T) {
	testCases := []TestCase{
		{
			name:      "Empty value",
			flag:      ScannerFlag,
			flagValue: "",
			want:      nil,
		},
		{
			name:      "String value",
			flag:      ScannerFlag,
			flagValue: "sast,secrets",
			want: []string{
				string(ScannerSAST),
				string(ScannerSecrets),
			},
		},
		{
			name: "Slice value",
			flag: ScannerFlag,
			flagValue: []string{
				"sast",
				"secrets",
			},
			want: []string{
				string(ScannerSAST),
				string(ScannerSecrets),
			},
		},
		{
			name: "Env value",
			flag: ScannerFlag,
			env: Env{
				key:   "BEARER_SCANNER",
				value: "sast,secrets",
			},
			want: []string{
				string(ScannerSAST),
				string(ScannerSecrets),
			},
		},
	}

	RunFlagTests(testCases, t)
}
