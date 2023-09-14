scopeCursor(req.params.oops)
scopeCursor(req.params.ok + x)
scopeCursor(x ? req.params.oops : y)
scopeCursor(req.params.ok ? x : y)

scopeNested(req.params.oops)
scopeResult(req.params.oops + x)
scopeNested(x ? req.params.oops : y)
scopeNested(req.params.oops ? x : y)

scopeResult(req.params.oops)
scopeResult(req.params.oops + x)
scopeResult(x ? req.params.oops : y)
scopeResult(req.params.ok ? x : y)
