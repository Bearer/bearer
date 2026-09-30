package orchestrator

import (
    "bytes"
    "context"
    "errors"
    "io"
    "os"
    "path/filepath"
    "strings"
    "testing"
    "time"

    "github.com/bearer/bearer/pkg/commands/process/filelist"
    "github.com/bearer/bearer/pkg/commands/process/filelist/files"
    "github.com/bearer/bearer/pkg/commands/process/orchestrator/work"
    "github.com/bearer/bearer/pkg/commands/process/settings"
    detectorapi "github.com/bearer/bearer/pkg/detectors"
    "github.com/bearer/bearer/pkg/report/coverage"
    reportdetectors "github.com/bearer/bearer/pkg/report/detectors"
    "github.com/bearer/bearer/pkg/report/writer"
    "github.com/hhatto/gocloc"
)

func TestCoverageReverifyMixedVisitationFixedPromise(t *testing.T) {
    root:=t.TempDir()
    body:=[]byte("package synthetic\n")
    fixed:=map[string]string{"input.go":coverage.Digest(body),"excluded.log":coverage.Digest(body)}
    for name:=range fixed{if err:=os.WriteFile(filepath.Join(root,name),body,0600);err!=nil{t.Fatal(err)}}
    cfg:=settings.Config{IgnoreGit:true,Worker:settings.WorkerOptions{TimeoutFileBytesPerSecond:10000,TimeoutFileMinimum:time.Second,TimeoutFileMaximum:time.Second}}
    cfg.Scan.Quiet=true
    list,err:=filelist.Discover(nil,root,&gocloc.Result{Files:map[string]*gocloc.ClocFile{}},cfg)
    if err!=nil||len(list.Files)!=2{t.Fatalf("fixture discovery: %+v %v",list,err)}
    calls:=0
    worker:=&fakeWorker{scan:func(req work.ProcessRequest)error{
        f,err:=os.Create(req.ReportPath);if err!=nil{return err}
        extractErr:=detectorapi.ExtractWithDetectors(context.Background(),req.Dir,req.File.FilePath,&writer.Detectors{File:f},nil,
            []detectorapi.InitializedDetector{{Type:reportdetectors.DetectorSimple,Detector:&reviewDetector{calls:&calls}}},nil,false,true)
        return errors.Join(extractErr,f.Close())
    }}
    o:=&Orchestrator{repository:work.Repository{Dir:root},config:&cfg,maxWorkersSemaphore:make(chan struct{},1),done:make(chan struct{}),pool:worker}
    output:=filepath.Join(t.TempDir(),"report.jsonl")
    err=o.Scan(output,list.Files)
    if calls!=1{t.Fatalf("actual detector visitation count=%d",calls)}
    if err==nil{t.Fatalf("mixed excluded work accepted; reconciliation=%v",reviewTerminal(t,output).ValidateInputs(fixed))}
    raw,readErr:=os.ReadFile(output);if readErr!=nil{t.Fatal(readErr)}
    if bytes.Contains(raw,[]byte(`"type":"processing_coverage"`)){t.Fatal("failed mixed work emitted a complete terminal record")}
    if len(fixed)!=2{t.Fatal("independent promise changed")}
    t.Logf("mixed run rejected, detector calls=%d fixed inputs=%d error=%v",calls,len(fixed),err)
}

type reverifyShortSink struct{dst io.Writer; writes int}
func(s *reverifyShortSink)Write(b []byte)(int,error){
    s.writes++
    if s.writes==1{return s.dst.Write(b[:len(b)-1])}
    return s.dst.Write(b)
}

func TestCoverageReverifyShortDiagnosticReachesTerminal(t *testing.T){
    var captured error
    var sink *reverifyShortSink
    o:=syntheticOrchestrator(t,func(req work.ProcessRequest)error{
        f,err:=os.Create(req.ReportPath);if err!=nil{return err}
        sink=&reverifyShortSink{dst:f}
        report:=&writer.Detectors{File:sink}
        captured=detectorapi.ExtractWithDetectors(context.Background(),req.Dir,req.File.FilePath,report,nil,
            []detectorapi.InitializedDetector{{Type:reportdetectors.DetectorSimple,Detector:coverageErrorDetector{}}},nil,false,true)
        report.AddError("input.go",errors.New("subsequent diagnostic"))
        if report.Err()==nil{return errors.New("test precondition: failure not retained")}
        return errors.Join(captured,report.Err(),f.Close())
    })
    out:=o.scanFile(io.Discard,files.File{FilePath:"input.go"})
    if sink==nil||sink.writes<2{t.Fatal("did not exercise short write followed by successful write")}
    if !errors.Is(captured,io.ErrShortWrite){t.Fatalf("real extraction did not retain short write: %v",captured)}
    if out.Status=="success"||!strings.Contains(out.Error,io.ErrShortWrite.Error()){t.Fatalf("short diagnostic failure disappeared: %+v",out)}
    t.Logf("real short diagnostic propagated to terminal: %+v",out)
}

