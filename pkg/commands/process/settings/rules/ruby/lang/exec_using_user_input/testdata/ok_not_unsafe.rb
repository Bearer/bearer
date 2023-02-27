Kernel.exec(y)

spawn(y)

IO.popen(y) {}

Process.exec(y)

Open3.popen3(["cmd", y], "abc") {}

Gem::Util.silent_system(x, y)

PTY.spawn("/bin/#{y}") {}

%x{/bin/#{y}}

`/bin/#{y}`


Shell.alias_command("foo", y) {}
Shell::CommandProcessor.alias_command(x, "/bin/#{y}") {}

Shell.def_system_command("foo", "bar", y) {}
Shell::CommandProcessor.def_system_command("foo", y) {}

shell = Shell.new(Dir.pwd)

processor1 = shell.command_processor
processor1.system(y)
