Dir["foo", base: request.env[:oops]]

Dir.chdir("/home/#{request.env[:oops]}")

File.exist?(request.env[:oops])

IO.readlines("/home/#{request.env[:oops]}")

Kernel.open(request.env[:oops], "w+") do
end

open(request.env[:oops])

PStore.new(request.env[:oops])

path = Pathname.new(request.env[:oops])
path + request.headers[:oops]
path / request.query_parameters[:oops]
path.join("a", request.body)
