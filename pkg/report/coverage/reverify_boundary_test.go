package coverage

import(
    "encoding/json"
    "strings"
    "testing"
)

func TestCoverageReverifyEveryRequiredNestedField(t *testing.T){
    r:=good()
    r.Discovery.Excluded=[]Exclusion{{Path:"omitted.go",Reason:"skip-path"}}
    raw,err:=json.Marshal(r);if err!=nil{t.Fatal(err)}
    levels:=[]struct{name string;keys []string;selectObject func(map[string]any)map[string]any}{
        {"root",[]string{"schema","run_id","configuration_sha256","discovery","outcomes","complete"},func(m map[string]any)map[string]any{return m}},
        {"discovery",[]string{"files","excluded","errors","complete"},func(m map[string]any)map[string]any{return m["discovery"].(map[string]any)}},
        {"exclusion",[]string{"path","reason"},func(m map[string]any)map[string]any{return m["discovery"].(map[string]any)["excluded"].([]any)[0].(map[string]any)}},
        {"outcome",[]string{"path","sha256","status","error"},func(m map[string]any)map[string]any{return m["outcomes"].([]any)[0].(map[string]any)}},
    }
    for _,level:=range levels{
        for _,key:=range level.keys{
            for _,mutation:=range []string{"missing","null","case_alias"}{
                t.Run(level.name+"/"+key+"/"+mutation,func(t *testing.T){
                    var m map[string]any
                    if err:=json.Unmarshal(raw,&m);err!=nil{t.Fatal(err)}
                    obj:=level.selectObject(m)
                    switch mutation{case "missing":delete(obj,key);case "null":obj[key]=nil;case "case_alias":obj[strings.ToUpper(key)]=obj[key]}
                    b,err:=json.Marshal(m);if err!=nil{t.Fatal(err)}
                    if _,err:=Decode(b);err==nil{t.Fatalf("invalid nested schema accepted: %s",b)}
                })
            }
        }
    }
    t.Run("escaped_duplicate_key",func(t *testing.T){
        mutated:=strings.Replace(string(raw),`"error":""`,`"error":"fault","\u0065rror":""`,1)
        if mutated==string(raw){t.Fatal("mutation missed outcome")}
        if _,err:=Decode([]byte(mutated));err==nil{t.Fatal("escaped duplicate outcome accepted")}
    })
    t.Run("fixed_promise_includes_exclusion",func(t *testing.T){
        decoded,err:=Decode(raw);if err!=nil{t.Fatal(err)}
        if decoded.ValidateInputs(map[string]string{"input.go":Digest([]byte("input")),"omitted.go":Digest(nil)})==nil{t.Fatal("excluded fixed promise shrank to discovery")}
    })
}
