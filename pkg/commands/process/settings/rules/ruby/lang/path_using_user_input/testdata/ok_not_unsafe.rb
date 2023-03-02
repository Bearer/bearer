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


Shell.cd(x)

Shell.default_system_path = x

shell = Shell.new(x, umask)

shell.pushdir(x, true)

processor1 = shell.command_processor
processor1.foreach(x, rs) {}

processor2 = Shell::CommandProcessor.new(shell)
processor2.test(:exists?, x, z)

processor2[:exists?, x, y, z]

processor.transact do
  test(:exists?, x)
end


render(partial: x, locals: { z: params[:ok] })
render_to_string({ file: "/templates/#{x}", locals: { z: params[:ok] } })

send_file x, type: "text/html"
