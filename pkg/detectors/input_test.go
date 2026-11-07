package detectors_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	detectorapi "github.com/bearer/bearer/pkg/detectors"
	"github.com/bearer/bearer/pkg/report/writer"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/bearer/bearer/pkg/commands/process/settings"
	"github.com/bearer/bearer/pkg/detectors/beego"
	"github.com/bearer/bearer/pkg/detectors/custom"
	gradle "github.com/bearer/bearer/pkg/detectors/dependencies/buildgradle/parser"
	"github.com/bearer/bearer/pkg/detectors/dependencies/depsbase"
	"github.com/bearer/bearer/pkg/detectors/dependencies/gemfile"
	"github.com/bearer/bearer/pkg/detectors/dependencies/gosum"
	"github.com/bearer/bearer/pkg/detectors/dependencies/ivy"
	paket "github.com/bearer/bearer/pkg/detectors/dependencies/paket-dependencies"
	"github.com/bearer/bearer/pkg/detectors/dependencies/requirements"
	"github.com/bearer/bearer/pkg/detectors/dependencies/yarnlock"
	"github.com/bearer/bearer/pkg/detectors/django"
	"github.com/bearer/bearer/pkg/detectors/internal/testhelper"
	"github.com/bearer/bearer/pkg/detectors/openapi"
	"github.com/bearer/bearer/pkg/detectors/openapi/queries"
	"github.com/bearer/bearer/pkg/detectors/php"
	"github.com/bearer/bearer/pkg/detectors/rails"
	"github.com/bearer/bearer/pkg/detectors/simple"
	"github.com/bearer/bearer/pkg/detectors/spring"
	"github.com/bearer/bearer/pkg/detectors/symfony"
	"github.com/bearer/bearer/pkg/detectors/types"
	"github.com/bearer/bearer/pkg/detectors/yamlconfig"
	"github.com/bearer/bearer/pkg/parser/nodeid"
	"github.com/bearer/bearer/pkg/util/file"
)

// TestCoverageF4ConsumerMatrix drives the public detector entrypoint with
// synthetic fixtures for each consumer moved to the bound input API. The
// fixtures intentionally contain ordinary URLs and package names only.
func TestCoverageF4ConsumerMatrix(t *testing.T) {
	tests := []struct {
		name      string
		filename  string
		contents  string
		companion map[string]string
	}{
		{
			name:     "custom-sql",
			filename: "schema.sql",
			contents: "CREATE TABLE synthetic_records (id integer);\n",
		},
		{
			name:     "simple",
			filename: "sample.go",
			contents: "package sample\nvar endpoint = \"https://example.test/api\"\n",
		},
		{
			name:     "openapi",
			filename: "api.yaml",
			contents: "openapi: 3.0.0\ninfo:\n  title: synthetic\n  version: 1.0.0\nservers:\n  - url: https://example.test\npaths: {}\n",
		},
		{
			name:     "php",
			filename: "sample.php",
			contents: "<?php\n$value = $_ENV['SYNTHETIC_URL'];\n",
		},
		{
			name:      "rails",
			filename:  filepath.Join("config", "database.yml"),
			contents:  "production:\n  adapter: postgresql\n",
			companion: map[string]string{"bin/rails": "#!/bin/sh\n"},
		},
		{
			name:      "spring-properties",
			filename:  "application.properties",
			contents:  "spring.datasource.url=jdbc:postgresql://example.test/db\n",
			companion: map[string]string{"pom.xml": "<project/>\n"},
		},
		{
			name:      "spring-yaml",
			filename:  "application.yml",
			contents:  "spring:\n  datasource:\n    url: jdbc:postgresql://example.test/db\n",
			companion: map[string]string{"pom.xml": "<project/>\n"},
		},
		{
			name:      "symfony",
			filename:  filepath.Join("config", "packages", "doctrine.yml"),
			contents:  "doctrine:\n  dbal:\n    driver: pdo_pgsql\n",
			companion: map[string]string{"composer.json": "{}\n"},
		},
		{
			name:     "yaml-config",
			filename: "config.yml",
			contents: "endpoint: https://example.test/config\n",
		},
		{
			name:     "requirements",
			filename: "requirements.txt",
			contents: "example-package==1.2.3\n",
		},
		{
			name:     "gemfile-lock",
			filename: "Gemfile.lock",
			contents: "GEM\n  remote: https://example.test\n  specs:\n    example-gem (1.2.3)\n\nBUNDLED WITH\n   2.5.0\n",
		},
		{
			name:     "gosum",
			filename: "go.sum",
			contents: "example.test/module v1.2.3 h1:synthetic\n",
		},
		{
			name:     "paket",
			filename: "paket.lock",
			contents: "NUGET\n  example-package (1.2.3)\n",
		},
		{
			name:     "yarn-lock",
			filename: "yarn.lock",
			contents: "# yarn lockfile v1\nexample-package@^1.2.3:\n  version \"1.2.3\"\n",
		},
		{
			name:     "ivy",
			filename: "ivy-report.xml",
			contents: "<ivy-report><dependencies><module organisation=\"example\" name=\"synthetic\"><revision name=\"1.2.3\"/></module></dependencies></ivy-report>\n",
		},
		{
			name:     "gradle",
			filename: "build.gradle",
			contents: "implementation 'example.group:synthetic:1.2.3'\n",
		},
	}

	tests = append(tests, tests[2])
	tests[len(tests)-1].name = "openapi-urls"
	tests = append(tests, tests[4])
	tests[len(tests)-1].name = "rails-storage"
	tests[len(tests)-1].filename = "config/storage.yml"
	tests[len(tests)-1].contents = "synthetic:\n  service: Disk\n  root: synthetic\n"
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			for _, mode := range []string{"valid", "leaf", "component"} {
				t.Run(mode, func(t *testing.T) {
					root := t.TempDir()
					writeF4Fixture(t, root, tt.filename, tt.contents)
					for name, contents := range tt.companion {
						writeF4Fixture(t, root, name, contents)
					}
					report := &testhelper.InMemoryReport{}
					visited := false
					err := file.IterateFilesListWithExclusions(root, []string{tt.filename}, false, true, func(*file.Path) (bool, error) { return true, nil }, func(info *file.FileInfo) error {
						visited = true
						if mode != "valid" {
							outside := t.TempDir()
							writeF4Fixture(t, outside, tt.filename, tt.contents)
							target, dest := info.AbsolutePath, filepath.Join(outside, tt.filename)
							if mode == "component" {
								target, dest = root, outside
							}
							if e := os.Rename(target, target+".held"); e != nil {
								return e
							}
							if e := os.Symlink(dest, target); e != nil {
								return e
							}
						}
						var detector types.Detector
						switch tt.name {
						case "custom-sql":
							d := custom.New(&nodeid.IntGenerator{}).(*custom.Detector)
							if e := d.CompileRules(map[string]*settings.Rule{"synthetic": {Languages: []string{"sql"}, DetectPresence: true, Patterns: []settings.RulePattern{{Pattern: "CREATE TABLE synthetic_records (id integer);"}}}}); e != nil {
								return e
							}
							detector = d
						case "simple":
							detector = simple.New()
						case "openapi":
							detector = openapi.New(&nodeid.IntGenerator{})
						case "openapi-urls":
							urls := queries.FindUrls(info)
							if mode == "valid" && len(urls) != 1 {
								return errors.New("valid URL consumer fixture not reached")
							}
							return nil
						case "php":
							detector = php.New(&nodeid.IntGenerator{})
						case "rails", "rails-storage":
							detector = rails.New(&nodeid.IntGenerator{})
						case "spring-properties", "spring-yaml":
							detector = spring.New()
						case "symfony":
							detector = symfony.New()
						case "yaml-config":
							detector = yamlconfig.New()
						default:
							funcs := map[string]func(*file.FileInfo) *depsbase.DiscoveredDependency{"requirements": requirements.Discover, "gemfile-lock": gemfile.Discover, "gosum": gosum.Discover, "paket": paket.Discover, "yarn-lock": yarnlock.Discover, "ivy": ivy.Discover, "gradle": gradle.Discover}
							result := funcs[tt.name](info)
							if mode == "valid" && result == nil {
								return errors.New("valid dependency consumer fixture failed")
							}
							return nil
						}
						_, e := detector.ProcessFile(info, &file.Path{AbsolutePath: root, RelativePath: ""}, report)
						return e
					}, nil)
					if !visited {
						t.Fatal("consumer never reached")
					}
					if mode == "valid" {
						if err != nil || len(report.Errors) != 0 {
							t.Fatalf("valid consumer failed: %v errors=%v", err, report.Errors)
						}
					} else if err == nil {
						t.Fatal("replacement consumer accepted")
					}
				})
			}
		})
	}
}

func writeF4Fixture(t *testing.T, root, name, contents string) {
	t.Helper()
	path := filepath.Join(root, name)
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(contents), 0o600); err != nil {
		t.Fatal(err)
	}
}

func TestCoverageContextExistsAcceptance(t *testing.T) {
	for _, family := range []string{"beego", "django", "spring"} {
		for _, mode := range []string{"present", "absent", "leaf-link", "component-link", "not-directory", "present-directory"} {
			if family == "django" && mode == "present-directory" {
				continue
			} // Django parses its marker; directories are not valid Python inputs.
			t.Run(family+"/"+mode, func(t *testing.T) {
				root := t.TempDir()
				os.WriteFile(filepath.Join(root, "input.go"), []byte("package synthetic\n"), 0600)
				project := filepath.Join(root, "project")
				os.Mkdir(project, 0700)
				var detector types.Detector
				var marker, body string
				switch family {
				case "beego":
					detector = beego.New()
					marker = "go.mod"
					body = "module example.test/synthetic\n"
				case "django":
					detector = django.New()
					marker = "manage.py"
					body = "import django\n"
				case "spring":
					detector = spring.New()
					marker = "pom.xml"
					body = "<project/>\n"
				}
				target := filepath.Join(project, marker)
				if mode != "absent" {
					os.WriteFile(target, []byte(body), 0600)
				}
				switch mode {
				case "leaf-link", "component-link":
					outside := t.TempDir()
					os.WriteFile(filepath.Join(outside, marker), []byte(body), 0600)
					if mode == "leaf-link" {
						os.Rename(target, target+".held")
						if err := os.Symlink(filepath.Join(outside, marker), target); err != nil {
							t.Fatal(err)
						}
					} else {
						os.Rename(project, project+".held")
						if err := os.Symlink(outside, project); err != nil {
							t.Fatal(err)
						}
					}
				case "not-directory":
					os.Remove(target)
					os.Remove(project)
					os.WriteFile(project, []byte("not a directory"), 0600)
				case "present-directory":
					os.Remove(target)
					os.Mkdir(target, 0700)
				}
				calls, visits, exclusions := 0, 0, 0
				var accepted bool
				var predicateError error
				err := file.IterateFilesListWithExclusions(root, []string{"input.go"}, false, true, func(dir *file.Path) (bool, error) {
					calls++
					accepted, predicateError = detector.AcceptDir(dir.Join("project"))
					return accepted, predicateError
				}, func(info *file.FileInfo) error { visits++; _, err := info.ReadBytes(); return err }, func(string) { exclusions++ })
				t.Logf("family=%s mode=%s calls=%d predicate=%v/%v visits=%d exclusions=%d terminal=%v", family, mode, calls, accepted, predicateError, visits, exclusions, err)
				if calls != 1 {
					t.Fatal("actual detector eligibility not invoked")
				}
				switch mode {
				case "present", "present-directory":
					if !accepted || err != nil || visits != 1 || exclusions != 0 {
						t.Fatal("valid context lost")
					}
				case "absent":
					if accepted || err != nil || visits != 0 || exclusions != 1 {
						t.Fatal("ordinary absence no longer compatible")
					}
				default:
					if accepted || err == nil || visits != 0 || exclusions != 0 {
						t.Fatal("unsafe metadata became successful processing or ordinary absence")
					}
				}
			})
		}
	}
}

func TestCoverageContextExistsTerminalMutation(t *testing.T) {
	for _, family := range []string{"beego", "django", "spring"} {
		for _, mode := range []string{"removed", "created", "leaf-replaced", "component-replaced"} {
			t.Run(family+"/"+mode, func(t *testing.T) {
				root := t.TempDir()
				os.WriteFile(filepath.Join(root, "input.go"), []byte("package synthetic\n"), 0600)
				project := filepath.Join(root, "project")
				os.Mkdir(project, 0700)
				var detector types.Detector
				var marker, body string
				switch family {
				case "beego":
					detector = beego.New()
					marker = "go.mod"
					body = "module example.test/synthetic\n"
				case "django":
					detector = django.New()
					marker = "manage.py"
					body = "import django\n"
				case "spring":
					detector = spring.New()
					marker = "pom.xml"
					body = "<project/>\n"
				}
				target := filepath.Join(project, marker)
				if mode != "created" {
					os.WriteFile(target, []byte(body), 0600)
				}
				called, visits := 0, 0
				err := file.IterateFilesListWithExclusions(root, []string{"input.go"}, false, true, func(dir *file.Path) (bool, error) {
					called++
					ok, e := detector.AcceptDir(dir.Join("project"))
					if e != nil {
						return false, e
					}
					if ok == (mode == "created") {
						t.Fatal("invalid initial marker state")
					}
					switch mode {
					case "removed":
						e = os.Remove(target)
					case "created":
						e = os.WriteFile(target, []byte(body), 0600)
					case "leaf-replaced", "component-replaced":
						outside := t.TempDir()
						os.WriteFile(filepath.Join(outside, marker), []byte(body), 0600)
						if mode == "leaf-replaced" {
							os.Rename(target, target+".held")
							e = os.Symlink(filepath.Join(outside, marker), target)
						} else {
							os.Rename(project, project+".held")
							e = os.Symlink(outside, project)
						}
					}
					return ok, e
				}, func(info *file.FileInfo) error { visits++; _, e := info.ReadBytes(); return e }, func(string) {})
				t.Logf("family=%s mode=%s called=%d visits=%d terminal=%v", family, mode, called, visits, err)
				if called != 1 || err == nil {
					t.Fatal("completed marker observation lost final reconciliation")
				}
				if mode == "created" && visits != 0 {
					t.Fatal("absent context counted as processing")
				}
			})
		}
	}
}

type contextHandoffBadSink struct{ short bool }

func (s contextHandoffBadSink) Write(b []byte) (int, error) {
	if s.short {
		return len(b) - 1, nil
	}
	return 0, errors.New("context sink failed")
}
func TestCoverageContextHandoffExtraction(t *testing.T) {
	for _, mode := range []string{"valid", "observation-free", "short", "write"} {
		t.Run(mode, func(t *testing.T) {
			root := t.TempDir()
			if err := os.WriteFile(filepath.Join(root, "input.go"), []byte("package synthetic\n"), 0600); err != nil {
				t.Fatal(err)
			}
			scope, err := file.OpenInputScope(root)
			if err != nil {
				t.Fatal(err)
			}
			defer scope.Close()
			bound, err := scope.Acquire("input.go", nil)
			if err != nil {
				t.Fatal(err)
			}
			var raw bytes.Buffer
			var sink io.Writer = &raw
			if mode == "short" || mode == "write" {
				sink = contextHandoffBadSink{mode == "short"}
			}
			scanners := []string{"sast"}
			if mode == "observation-free" {
				scanners = nil
			}
			err = detectorapi.ExtractWithInput(context.Background(), root, "input.go", &writer.Detectors{File: sink}, nil, scanners, nil, false, true, bound.Expectation())
			if mode == "short" || mode == "write" {
				if err == nil {
					t.Fatal("diagnostic failure accepted")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			dec := json.NewDecoder(&raw)
			count := 0
			for {
				var record json.RawMessage
				e := dec.Decode(&record)
				if e == io.EOF {
					break
				}
				if e != nil {
					t.Fatal(e)
				}
				var header struct{ Type string }
				if e = json.Unmarshal(record, &header); e != nil {
					t.Fatal(e)
				}
				if header.Type == file.ContextReceiptType {
					receipt, e := file.DecodeContextReceipt(record)
					if e != nil {
						t.Fatal(e)
					}
					if mode == "observation-free" && len(receipt.Observations) != 0 {
						t.Fatal("invented observation")
					}
					count++
				}
			}
			if count != 1 {
				t.Fatalf("receipt count=%d", count)
			}
		})
	}
}
