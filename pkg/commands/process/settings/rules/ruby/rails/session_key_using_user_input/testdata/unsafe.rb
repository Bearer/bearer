session[params[:key]] = 42

session[request.env[:key]] = 42

session["test-#{cookies["oops"]}"] = 42
