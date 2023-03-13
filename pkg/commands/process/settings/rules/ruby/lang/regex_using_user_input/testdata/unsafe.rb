/#{params[:oops]}.*/

%r{abc#{params[:oops]}def}

Regexp.new(params[:oops])
