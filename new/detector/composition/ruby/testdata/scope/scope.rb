scope_cursor(params[:oops])
scope_cursor(x ? params[:ok] : y)
scope_cursor(params[:ok] ? x : y)

scope_nested(params[:oops])
scope_nested(x ? params[:oops] : y)
scope_nested(params[:oops] ? x : y)

scope_result(params[:oops])
scope_result(x ? params[:oops] : y)
scope_result(params[:ok] ? x : y)
