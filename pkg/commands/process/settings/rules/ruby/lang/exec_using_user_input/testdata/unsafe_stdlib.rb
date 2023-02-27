# stdlib

Kernel.exec(params[:oops])

spawn(params[:oops])

IO.popen(params[:oops]) {}

Process.exec(params[:oops])

Open3.popen3(["cmd", params[:oops]], "abc") {}

Gem::Util.silent_system(x, params[:oops])

PTY.spawn("/bin/#{params[:oops]}") {}

%x{/bin/#{params[:oops]}}

`/bin/#{params[:oops]}`
