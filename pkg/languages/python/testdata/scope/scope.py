scopeCursor(request.GET.get('oops'))
scopeCursor(x + request.GET.get("ok"))
scopeCursor(request.GET.get('oops') if x else y)
scopeCursor(x if request.GET.get('ok') else y)
scopeCursor(request.GET.get('oops') or y) # wrong

scopeNested(request.GET.get('oops'))
scopeNested(x + request.GET.get('oops'))
scopeNested(request.GET.get('oops') if x else y)
scopeNested(x if request.GET.get('oops') else y)
scopeNested(request.GET.get('oops') or y)

scopeResult(request.GET.get('oops'))
scopeResult(x + request.GET.get('oops'))
scopeResult(request.GET.get('oops') if x else y)
scopeResult(x if request.GET.get('ok') else y)
scopeResult(request.GET.get('oops') or y)