package coverage

import (
    "encoding/json"
    "strings"
    "testing"
)

func TestCoverageReviewRejectsPartialAndAliasedOutcomes(t *testing.T) {
    record := Record{
        Schema: Schema, RunID: Digest([]byte("review-run")), ConfigurationSHA256: Digest([]byte("review-config")),
        Discovery: Discovery{Files: []string{"input.go"}, Excluded: []Exclusion{}, Errors: []string{}, Complete: true},
        Outcomes: []Outcome{{Path: "input.go", SHA256: Digest([]byte("fixed input")), Status: "success", Error: ""}}, Complete: true,
    }
    encoded, err := json.Marshal(record)
    if err != nil { t.Fatal(err) }
    fixedPromise := map[string]string{"input.go": Digest([]byte("fixed input"))}
    cases := []struct { name, from, to string }{
        {"missing_error", `,"error":""`, ""},
        {"null_error", `"error":""`, `"error":null`},
        {"case_alias_overwrites_error", `"error":""`, `"error":"actual failure","ERROR":""`},
        {"case_alias_overwrites_status", `"status":"success"`, `"status":"timeout","STATUS":"success"`},
    }
    for _, tc := range cases {
        t.Run(tc.name, func(t *testing.T) {
            raw := strings.Replace(string(encoded), tc.from, tc.to, 1)
            if raw == string(encoded) { t.Fatal("probe mutation did not apply") }
            got, err := Decode([]byte(raw))
            if err == nil { err = got.ValidateInputs(fixedPromise) }
            if err == nil { t.Errorf("malformed/partial outcome accepted against fixed input promise: %s", raw) }
        })
    }
}

