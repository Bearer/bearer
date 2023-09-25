scopeCursor(request.getParameter("oops"));
scopeCursor(x + request.getParameter("ok"));
scopeCursor(x ? request.getParameter("oops") : y);
scopeCursor(request.getParameter("ok") ? x : y);

scopeNested(request.getParameter("oops"));
scopeNested(x + request.getParameter("oops"));
scopeNested(x ? request.getParameter("oops") : y);
scopeNested(request.getParameter("oops") ? x : y);

scopeResult(request.getParameter("oops"));
scopeResult(x + request.getParameter("oops"));
scopeResult(x ? request.getParameter("oops") : y);
scopeResult(request.getParameter("ok") ? x : y);
