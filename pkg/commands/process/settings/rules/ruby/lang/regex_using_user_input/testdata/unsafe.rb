/#{params[:oops]}.*/

%r{abc#{params[:oops]}def}

Regex.new(params[:oops])
