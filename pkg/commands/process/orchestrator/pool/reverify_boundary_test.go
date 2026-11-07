package pool

import(
    "context"
    "encoding/json"
    "errors"
    "net/http"
    "os"
    "path/filepath"
    "testing"
    "time"

    "github.com/bearer/bearer/pkg/commands/process/filelist/files"
    "github.com/bearer/bearer/pkg/commands/process/orchestrator/work"
    "github.com/bearer/bearer/pkg/commands/process/orchestrator/worker"
    "github.com/bearer/bearer/pkg/scanner/stats"
)

func TestCoverageReverifyWorkerTimeoutClassification(t *testing.T){
    root:=t.TempDir()
    if err:=os.WriteFile(filepath.Join(root,"input.go"),[]byte("package synthetic\n"),0600);err!=nil{t.Fatal(err)}
    ctx,cancel:=context.WithDeadline(context.Background(),time.Now().Add(-time.Second));defer cancel()
    w:=&worker.Worker{}
    fileStats,workerErr:=w.Scan(ctx,work.ProcessRequest{Repository:work.Repository{Dir:root},File:files.File{FilePath:"input.go"},ReportPath:filepath.Join(t.TempDir(),"report.jsonl")})
    if !errors.Is(workerErr,worker.ErrorTimeoutReached){t.Fatalf("actual worker did not report its timeout sentinel: %v",workerErr)}
    payload,err:=json.Marshal(work.ProcessResponse{FileStats:fileStats,Error:workerErr.Error()});if err!=nil{t.Fatal(err)}
    body,response:=responseBody(string(payload))
    process:=syntheticProcess(coverageRoundTrip(func(*http.Request)(*http.Response,error){return response,nil}))
    p:=&Pool{available:[]*Process{process}}
    defer p.Close()
    got:=p.Scan(processRequest())
    waitClosed(t,body.closed)
    if got==nil||got.Error()!=workerErr.Error(){t.Fatalf("actual worker error changed in transport: %v want=%v",got,workerErr)}
    reason:=translateErrorForStats(got.Error(),false)
    if reason!=stats.TimeoutReason{
        t.Errorf("actual worker timeout misclassified as %q; returned error=%q",reason.String(),got.Error())
    }
}

func TestCoverageReverifyResponseExactFields(t *testing.T){
    for _,raw:=range []string{
        `{"filestats":null,"Error":""}`,
        `{"FileStats":null,"error":""}`,
        `{"FileStats":[],"Error":""}`,
        `{"FileStats":null,"Error":null}`,
        `{"FileStats":null,"Error":false}`,
        `{"FileStats":null,"Error":"","error":""}`,
    }{
        t.Run(raw,func(t *testing.T){
            body,response:=responseBody(raw)
            p:=syntheticProcess(coverageRoundTrip(func(*http.Request)(*http.Response,error){return response,nil}))
            if got,err:=p.Scan(processRequest());got!=nil||err==nil{t.Fatalf("malformed worker fields accepted: %#v %v",got,err)}
            waitClosed(t,body.closed)
        })
    }
}

