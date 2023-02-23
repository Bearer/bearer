Dir["foo", x]

Dir.chdir("/home/#{x}")

event = not_from_handler
File.exist?(event["oops"])

IO.readlines("/home/#{x}")

Kernel.open(x, "w+") do
end

open(x)

PStore.new(x)

path = Pathname.new(x)
path + x
path / x
path.join("a", x)


Rails.root.join(x)

render(partial: x, locals: { z: params[:ok] })
render_to_string({ file: "/templates/#{x}", locals: { z: params[:ok] } })
