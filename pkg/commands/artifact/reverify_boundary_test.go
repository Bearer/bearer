package artifact

import (
    "bytes"
    "encoding/json"
    "os"
    "testing"

    "github.com/bearer/bearer/pkg/commands/process/filelist/files"
    "github.com/bearer/bearer/pkg/commands/process/settings"
    "github.com/bearer/bearer/pkg/flag"
    "github.com/bearer/bearer/pkg/report/coverage"
    "github.com/bearer/bearer/pkg/report/detections"
)

func TestCoverageReverifyCacheProducerTypes(t *testing.T) {
    cases := []struct{name, record string}{
        {"detector_language_array", `{"type":"schema","detector_type":"simple","detector_language":[],"source":{"filename":"input.go"},"value":{}}`},
        {"source_language_boolean", `{"type":"schema","detector_type":"simple","source":{"filename":"input.go","language":false},"value":{}}`},
    }
    for _, mode := range []string{flag.ReportSecurity, flag.ReportDetectors} {
        for _, tc := range cases {
            t.Run(mode+"/"+tc.name, func(t *testing.T) {
                var typed detections.Detection
                if err := json.Unmarshal([]byte(tc.record), &typed); err == nil {
                    t.Fatal("negative fixture unexpectedly satisfies the existing producer envelope")
                } else { t.Logf("producer envelope rejects fixture: %v", err) }
                r := coverageCachedRunnerConfigured(t, `{"type":"file_list","filenames":["input.go"]}`+"\n"+tc.record+"\n", false, func(c *settings.Config) {
                    c.Report.Report = mode
                    if mode == flag.ReportDetectors { c.Report.Format = flag.FormatJSON }
                }).(*runner)
                if !r.CacheUsed() { t.Fatal("completed cache not actually selected") }
                failed, err := r.Report(nil, nil)
                if err == nil {
                    delivered, readErr := os.ReadFile(r.scanSettings.Report.Output)
                    t.Errorf("malformed typed legacy record accepted: failed=%v delivered=%s readErr=%v", failed, delivered, readErr)
                }
            })
        }
    }
    t.Run("valid_producer_control", func(t *testing.T) {
        r := coverageCachedRunner(t, `{"type":"file_list","filenames":["input.go"]}
{"type":"schema","detector_type":"simple","detector_language":"go","source":{"filename":"input.go","language":"Go"},"value":{}}
`, false).(*runner)
        if _, err := r.Report(nil,nil); err != nil { t.Fatal(err) }
        b, err := os.ReadFile(r.scanSettings.Report.Output)
        if err != nil || !bytes.Contains(b, []byte(`"coverage":null`)) { t.Fatalf("legacy control: %s %v",b,err) }
    })
}

func TestCoverageReverifyCacheUnicodeAlias(t *testing.T) {
    raw := `{"type":"file_list","filenames":["input.go"],"filenameſ":[]}`
    var decoded detections.FileListDetection
    if err := json.Unmarshal([]byte(raw), &decoded); err != nil { t.Fatal(err) }
    if len(decoded.Filenames) != 0 { t.Fatal("fixture did not exercise encoding/json Unicode field folding") }
    r := coverageCachedRunner(t, raw+"\n", false).(*runner)
    if !r.CacheUsed() { t.Fatal("cache not selected") }
    if failed, err := r.Report(nil,nil); err == nil {
        b, readErr := os.ReadFile(r.scanSettings.Report.Output)
        t.Errorf("Unicode alias overwrote validated file list and was delivered: failed=%v output=%s readErr=%v", failed,b,readErr)
    }
}

func TestCoverageReverifyFreshSchemaAtDelivery(t *testing.T) {
    fixedDigest := coverage.Digest([]byte("independently fixed synthetic input"))
    d := coverage.Discovery{Files:[]string{"input.go"}, Excluded:[]coverage.Exclusion{}, Errors:[]string{},Complete:true}
    record := coverage.Record{Schema:coverage.Schema, RunID:coverage.Digest([]byte("new synthetic run")),ConfigurationSHA256:coverage.Digest(nil),Discovery:d,Outcomes:[]coverage.Outcome{{Path:"input.go",SHA256:fixedDigest,Status:"success",Error:""}},Complete:true}
    b, err := json.Marshal(record)
    if err != nil { t.Fatal(err) }
    cases:=[]struct{name string; mutate func(map[string]any)}{
        {"control",func(map[string]any){}},
        {"nested_exclusion_null_reason",func(m map[string]any){m["discovery"].(map[string]any)["excluded"]=[]any{map[string]any{"path":"skip.go","reason":nil}}}},
        {"nested_discovery_alias",func(m map[string]any){m["discovery"].(map[string]any)["Complete"]=true}},
        {"nested_outcome_unknown",func(m map[string]any){m["outcomes"].([]any)[0].(map[string]any)["ignored_fault"]="failure"}},
        {"duplicate_terminal",func(m map[string]any){m["outcomes"]=append(m["outcomes"].([]any),m["outcomes"].([]any)[0])}},
    }
    for _, tc := range cases {
        t.Run(tc.name,func(t *testing.T){
            var m map[string]any
            if err:=json.Unmarshal(b,&m);err!=nil{t.Fatal(err)}
            tc.mutate(m)
            payload,err:=json.Marshal(map[string]any{"type":"processing_coverage","coverage":m})
            if err!=nil{t.Fatal(err)}
            r:=coverageCachedRunner(t,`{"type":"file_list","filenames":["input.go"]}`+"\n"+string(payload)+"\n",false).(*runner)
            r.reuseDetection=false
            r.discovery=d
            _,err=r.Report([]files.File{{FilePath:"input.go"}},nil)
            if tc.name!="control" {
                if err==nil{t.Fatal("malformed fresh evidence reached delivery")}
                if _,statErr:=os.Stat(r.scanSettings.Report.Output);!os.IsNotExist(statErr){t.Fatalf("rejected evidence was delivered: %v",statErr)}
                return
            }
            if err!=nil{t.Fatal(err)}
            delivered,err:=os.ReadFile(r.scanSettings.Report.Output);if err!=nil{t.Fatal(err)}
            var envelope struct{Coverage *coverage.Record `json:"coverage"`}
            if err=json.Unmarshal(delivered,&envelope);err!=nil{t.Fatal(err)}
            if err=envelope.Coverage.ValidateInputs(map[string]string{"input.go":fixedDigest});err!=nil{t.Fatal(err)}
        })
    }
}

