package foo

import (
	"net/http"
)

func scopeCursor(s any) any { return s } // nolint: unused
func scopeNested(s any) any { return s } // nolint: unused
func scopeResult(s any) any { return s } // nolint: unused

func loginHandler(w http.ResponseWriter, request *http.Request) { // nolint: unused
	var x string
	y := map[string]string{}

	scopeCursor(request.FormValue("oops"))   // nolint: staticcheck
	scopeCursor(x + request.FormValue("ok")) // nolint: staticcheck

	scopeNested(request.FormValue("oops"))     // nolint: staticcheck
	scopeNested(x + request.FormValue("oops")) // nolint: staticcheck
	scopeNested(y[request.FormValue("oops")])  // nolint: staticcheck

	scopeResult(request.FormValue("oops"))     // nolint: staticcheck
	scopeResult(x + request.FormValue("oops")) // nolint: staticcheck
	scopeResult(y[request.FormValue("ok")])    // nolint: staticcheck
}

func main() { // nolint: unused
	var req *http.Request
	var x string
	y := map[string]string{}

	scopeCursor(req.FormValue("oops"))   // nolint: staticcheck
	scopeCursor(x + req.FormValue("ok")) // nolint: staticcheck

	scopeNested(req.FormValue("oops"))     // nolint: staticcheck
	scopeNested(x + req.FormValue("oops")) // nolint: staticcheck
	scopeNested(y[req.FormValue("oops")])  // nolint: staticcheck

	scopeResult(req.FormValue("oops"))     // nolint: staticcheck
	scopeResult(x + req.FormValue("oops")) // nolint: staticcheck
	scopeResult(y[req.FormValue("ok")])    // nolint: staticcheck
}
