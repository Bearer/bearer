Dir["foo", base: params[:oops]]

Dir.chdir("/home/#{params[:oops]}")

File.exist?(params[:oops])

IO.readlines("/home/#{params[:oops]}")

Kernel.open(params[:oops], "w+") do
end

open(params[:oops])

PStore.new(params[:oops])

path = Pathname.new(params[:oops])
path + params[:two]
path / params[:three]
path.join("a", params[:four])

Rails.root.join(params[:oops])

Gem::Util.traverse_parents(params[:oops]) {}
