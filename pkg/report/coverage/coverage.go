// Package coverage accounts for actual same-run terminal processing. A consumer
// must supply its independently fixed promise; discovery is never that promise.
package coverage

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"path"
	"strings"
)

const Schema = "bearer-processing-v1"

type Exclusion struct {
	// Repository scope is reserved for inventory outside the fixed target.
	Scope  string `json:"scope,omitempty"`
	Path   string `json:"path"`
	Reason string `json:"reason"`
}
type Discovery struct {
	Files    []string    `json:"files"`
	Excluded []Exclusion `json:"excluded"`
	Errors   []string    `json:"errors"`
	Complete bool        `json:"complete"`
}
type Outcome struct {
	Path   string `json:"path"`
	SHA256 string `json:"sha256"`
	Status string `json:"status"`
	Error  string `json:"error"`
}
type Record struct {
	Schema              string    `json:"schema"`
	RunID               string    `json:"run_id"`
	ConfigurationSHA256 string    `json:"configuration_sha256"`
	Discovery           Discovery `json:"discovery"`
	Outcomes            []Outcome `json:"outcomes"`
	Complete            bool      `json:"complete"`
}

func Digest(data []byte) string { h := sha256.Sum256(data); return hex.EncodeToString(h[:]) }
func validDigest(s string) bool {
	b, e := hex.DecodeString(s)
	return e == nil && len(b) == sha256.Size && s == strings.ToLower(s)
}
func validPath(s string) bool {
	return s != "" && !strings.Contains(s, "\\") && !strings.ContainsRune(s, 0) && !path.IsAbs(s) && path.Clean(s) == s && s != ".." && !strings.HasPrefix(s, "../")
}

// Validate never substitutes discovered files when the caller's promise is absent.
// Complete describes internal execution; coverage of a caller's promise requires
// this reconciliation as well as successful delivery and an accepted process exit.
func (r *Record) Validate(promised []string) error {
	if r == nil || r.Schema != Schema || !validDigest(r.RunID) || !validDigest(r.ConfigurationSHA256) || !r.Complete {
		return errors.New("coverage terminal record unavailable or incomplete")
	}
	d := r.Discovery
	if !d.Complete || d.Files == nil || d.Excluded == nil || d.Errors == nil || len(d.Errors) > 0 || r.Outcomes == nil || promised == nil {
		return errors.New("coverage discovery or promise unavailable/incomplete")
	}
	discovered := map[string]bool{}
	for _, p := range d.Files {
		if !validPath(p) || discovered[p] {
			return fmt.Errorf("invalid/duplicate discovered path: %q", p)
		}
		discovered[p] = true
	}
	excluded := map[string]bool{}
	for _, e := range d.Excluded {
		if e.Scope != "" && e.Scope != "repository" {
			return errors.New("invalid exclusion scope")
		}
		if e.Scope == "repository" && e.Reason != "outside-target" {
			return errors.New("invalid outside-target exclusion")
		}
		key := e.Scope + "\x00" + e.Path
		if !validPath(e.Path) || e.Reason == "" || excluded[key] || (e.Scope == "" && discovered[e.Path]) {
			return errors.New("invalid/duplicate exclusion")
		}
		excluded[key] = true
	}
	terminal := map[string]bool{}
	for _, o := range r.Outcomes {
		if !discovered[o.Path] || terminal[o.Path] {
			return errors.New("unknown/duplicate terminal outcome")
		}
		if o.Status != "success" || o.Error != "" || !validDigest(o.SHA256) {
			return errors.New("unsuccessful terminal outcome")
		}
		terminal[o.Path] = true
	}
	if len(terminal) != len(discovered) {
		return errors.New("missing terminal outcome")
	}
	seen := map[string]bool{}
	for _, p := range promised {
		if !validPath(p) || seen[p] {
			return errors.New("invalid/duplicate promised input")
		}
		seen[p] = true
		if !terminal[p] {
			return fmt.Errorf("promised input not processed: %q", p)
		}
	}
	return nil
}

// Decode checks the record shape. Semantic validation still needs the promise.
func Decode(data []byte) (*Record, error) {
	if err := CheckJSON(data); err != nil {
		return nil, err
	}
	if err := checkRecordShape(data); err != nil {
		return nil, err
	}
	var r Record
	d := json.NewDecoder(bytes.NewReader(data))
	d.DisallowUnknownFields()
	if err := d.Decode(&r); err != nil {
		return nil, err
	}
	var extra any
	if err := d.Decode(&extra); err != io.EOF {
		return nil, errors.New("trailing coverage data")
	}
	return &r, nil
}

// WriteFull converts a short nil-error write into a failure.
func WriteFull(w io.Writer, data []byte) error {
	n, err := w.Write(data)
	if err != nil {
		return err
	}
	if n != len(data) {
		return io.ErrShortWrite
	}
	return nil
}

// ValidateInputs also binds outcomes to the independently fixed input bytes.
func (r *Record) ValidateInputs(promised map[string]string) error {
	if promised == nil {
		return errors.New("input promise unavailable")
	}
	paths := make([]string, 0, len(promised))
	for p, d := range promised {
		if !validDigest(d) {
			return errors.New("invalid promised digest")
		}
		paths = append(paths, p)
	}
	if err := r.Validate(paths); err != nil {
		return err
	}
	for _, o := range r.Outcomes {
		if digest, ok := promised[o.Path]; ok && digest != o.SHA256 {
			return errors.New("promised input bytes differ")
		}
	}
	return nil
}

// CheckJSON rejects duplicate keys as well as malformed/trailing JSON. A map
// decoder alone could silently replace a failed outcome with a later success.
func CheckJSON(data []byte) error {
	d := json.NewDecoder(bytes.NewReader(data))
	d.UseNumber()
	var value func() error
	value = func() error {
		token, err := d.Token()
		if err != nil {
			return err
		}
		delim, ok := token.(json.Delim)
		if !ok {
			return nil
		}
		switch delim {
		case '{':
			keys := map[string]bool{}
			for d.More() {
				key, err := d.Token()
				if err != nil {
					return err
				}
				s, ok := key.(string)
				if !ok || keys[s] {
					return errors.New("duplicate/invalid JSON key")
				}
				keys[s] = true
				if err = value(); err != nil {
					return err
				}
			}
			end, err := d.Token()
			if err != nil {
				return err
			}
			if end != json.Delim('}') {
				return errors.New("incomplete JSON object")
			}
		case '[':
			for d.More() {
				if err = value(); err != nil {
					return err
				}
			}
			end, err := d.Token()
			if err != nil {
				return err
			}
			if end != json.Delim(']') {
				return errors.New("incomplete JSON array")
			}
		default:
			return errors.New("unexpected JSON delimiter")
		}
		return nil
	}
	if err := value(); err != nil {
		return err
	}
	if _, err := d.Token(); err != io.EOF {
		return errors.New("trailing JSON data")
	}
	return nil
}

// Shape validation precedes struct decoding: encoding/json otherwise accepts
// case-insensitive field aliases and converts null/missing scalar fields to zero.
func shapeObject(value any, required map[string]string, optional map[string]string) (map[string]any, error) {
	m, ok := value.(map[string]any)
	if !ok || m == nil {
		return nil, errors.New("coverage object required")
	}
	for key := range required {
		if _, ok := m[key]; !ok {
			return nil, fmt.Errorf("missing coverage field %s", key)
		}
	}
	for key, v := range m {
		kind, ok := required[key]
		if !ok {
			kind, ok = optional[key]
		}
		if !ok || v == nil {
			return nil, fmt.Errorf("unknown/null coverage field %s", key)
		}
		valid := false
		switch kind {
		case "string":
			_, valid = v.(string)
		case "bool":
			_, valid = v.(bool)
		case "array":
			_, valid = v.([]any)
		case "object":
			_, valid = v.(map[string]any)
		}
		if !valid {
			return nil, fmt.Errorf("invalid coverage field type %s", key)
		}
	}
	return m, nil
}
func checkRecordShape(data []byte) error {
	var value any
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	root, err := shapeObject(value, map[string]string{"schema": "string", "run_id": "string", "configuration_sha256": "string", "discovery": "object", "outcomes": "array", "complete": "bool"}, nil)
	if err != nil {
		return err
	}
	discovery, err := shapeObject(root["discovery"], map[string]string{"files": "array", "excluded": "array", "errors": "array", "complete": "bool"}, nil)
	if err != nil {
		return err
	}
	for _, key := range []string{"files", "errors"} {
		for _, v := range discovery[key].([]any) {
			if _, ok := v.(string); !ok {
				return errors.New("coverage string entry required")
			}
		}
	}
	for _, v := range discovery["excluded"].([]any) {
		if _, err := shapeObject(v, map[string]string{"path": "string", "reason": "string"}, map[string]string{"scope": "string"}); err != nil {
			return err
		}
	}
	for _, v := range root["outcomes"].([]any) {
		if _, err := shapeObject(v, map[string]string{"path": "string", "sha256": "string", "status": "string", "error": "string"}, nil); err != nil {
			return err
		}
	}
	return nil
}
