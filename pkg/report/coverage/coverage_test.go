package coverage

import (
	"encoding/json"
	"errors"
	"testing"
)

func good() Record {
	return Record{Schema: Schema, RunID: Digest([]byte("run")), ConfigurationSHA256: Digest([]byte("config")), Discovery: Discovery{Files: []string{"input.go"}, Excluded: []Exclusion{}, Errors: []string{}, Complete: true}, Outcomes: []Outcome{{Path: "input.go", SHA256: Digest([]byte("input")), Status: "success"}}, Complete: true}
}
func TestCoverageReconciliation(t *testing.T) {
	cases := []struct {
		name    string
		mutate  func(*Record)
		promise []string
		want    bool
	}{
		{"complete", func(r *Record) {}, []string{"input.go"}, true},
		{"omitted-promise", func(r *Record) {}, []string{"input.go", "omitted.go"}, false},
		{"exclusion-not-processing", func(r *Record) { r.Discovery.Excluded = []Exclusion{{Path: "excluded.go", Reason: "minified-js"}} }, []string{"input.go", "excluded.go"}, false},
		{"missing-outcome", func(r *Record) { r.Outcomes = []Outcome{} }, []string{"input.go"}, false},
		{"duplicate-outcome", func(r *Record) { r.Outcomes = append(r.Outcomes, r.Outcomes[0]) }, []string{"input.go"}, false},
		{"unknown-outcome", func(r *Record) { r.Outcomes[0].Path = "unknown.go" }, []string{"input.go"}, false},
		{"error", func(r *Record) { r.Outcomes[0].Error = "synthetic failure" }, []string{"input.go"}, false},
		{"timeout", func(r *Record) { r.Outcomes[0].Status = "timeout" }, []string{"input.go"}, false},
		{"cancelled", func(r *Record) { r.Complete = false }, []string{"input.go"}, false},
		{"discovery-error", func(r *Record) { r.Discovery.Errors = []string{"stat error"} }, []string{"input.go"}, false},
		{"missing-discovery", func(r *Record) { r.Discovery = Discovery{} }, []string{"input.go"}, false},
		{"duplicate-discovery", func(r *Record) { r.Discovery.Files = append(r.Discovery.Files, "input.go") }, []string{"input.go"}, false},
		{"no-promise", func(r *Record) {}, nil, false},
		{"invalid-path", func(r *Record) { r.Discovery.Files[0] = "../escape" }, []string{"input.go"}, false},
		{"missing-digest", func(r *Record) { r.Outcomes[0].SHA256 = "" }, []string{"input.go"}, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			r := good()
			c.mutate(&r)
			if got := r.Validate(c.promise) == nil; got != c.want {
				t.Fatalf("accepted=%v want=%v", got, c.want)
			}
		})
	}
}
func TestCoverageEmptyAndInputBinding(t *testing.T) {
	r := good()
	if err := r.ValidateInputs(map[string]string{"input.go": Digest([]byte("input"))}); err != nil {
		t.Fatal(err)
	}
	if r.ValidateInputs(map[string]string{"input.go": Digest([]byte("changed"))}) == nil {
		t.Fatal("changed bytes accepted")
	}
	r.Discovery.Files = []string{}
	r.Outcomes = []Outcome{}
	if err := r.Validate([]string{}); err != nil {
		t.Fatal(err)
	}
	r.Discovery.Files = nil
	if r.Validate([]string{}) == nil {
		t.Fatal("missing discovery accepted as empty")
	}
}
func TestCoverageDecodeAndWrite(t *testing.T) {
	b, _ := json.Marshal(good())
	if _, err := Decode(b); err != nil {
		t.Fatal(err)
	}
	for _, bad := range [][]byte{b[:len(b)-1], append(append([]byte{}, b...), b...), []byte(`{"unknown":1}`), []byte(`{"complete":false,"complete":true}`)} {
		if _, err := Decode(bad); err == nil {
			t.Fatal("malformed record accepted")
		}
	}
	for _, w := range []badWriter{{false}, {true}} {
		if WriteFull(w, []byte("bytes")) == nil {
			t.Fatal("failed write accepted")
		}
	}
}

type badWriter struct{ fault bool }

func (w badWriter) Write(b []byte) (int, error) {
	if w.fault {
		return 0, errors.New("synthetic write error")
	}
	return len(b) - 1, nil
}

func TestCoverageExactRequiredSchema(t *testing.T) {
	b, _ := json.Marshal(good())
	for _, test := range []struct {
		name   string
		mutate func(map[string]any)
	}{
		{"missing complete", func(m map[string]any) { delete(m, "complete") }},
		{"null complete", func(m map[string]any) { m["complete"] = nil }},
		{"wrong boolean", func(m map[string]any) { m["complete"] = "true" }},
		{"unknown alias", func(m map[string]any) { m["Complete"] = true }},
		{"null files", func(m map[string]any) { m["discovery"].(map[string]any)["files"] = nil }},
		{"nonstring path", func(m map[string]any) { m["discovery"].(map[string]any)["files"] = []any{3} }},
		{"missing outcome error", func(m map[string]any) { delete(m["outcomes"].([]any)[0].(map[string]any), "error") }},
		{"null outcome error", func(m map[string]any) { m["outcomes"].([]any)[0].(map[string]any)["error"] = nil }},
		{"wrong outcome status", func(m map[string]any) { m["outcomes"].([]any)[0].(map[string]any)["status"] = true }},
	} {
		t.Run(test.name, func(t *testing.T) {
			var m map[string]any
			_ = json.Unmarshal(b, &m)
			test.mutate(m)
			bad, _ := json.Marshal(m)
			if _, err := Decode(bad); err == nil {
				t.Fatal("invalid shape accepted")
			}
		})
	}
	if _, err := Decode(b); err != nil {
		t.Fatal(err)
	}
}
func TestCoverageRootFileAndUnsafePaths(t *testing.T) {
	for _, p := range []string{".", "../outside", "/absolute", "a/../b", "", `a\b`} {
		record := good()
		record.Discovery.Files = []string{p}
		record.Outcomes[0].Path = p
		err := record.ValidateInputs(map[string]string{p: Digest([]byte("input"))})
		if (err == nil) != (p == ".") {
			t.Fatalf("path %q error %v", p, err)
		}
	}
	record := good()
	record.Discovery.Excluded = []Exclusion{{Scope: "repository", Path: "input.go", Reason: "outside-target"}}
	if err := record.Validate([]string{"input.go"}); err != nil {
		t.Fatal(err)
	}
	if err := record.Validate([]string{"../input.go"}); err == nil {
		t.Fatal("outside-target promise accepted")
	}
	record.Discovery.Excluded[0].Reason = "minified"
	if err := record.Validate([]string{"input.go"}); err == nil {
		t.Fatal("scope used to hide in-target exclusion")
	}
}
