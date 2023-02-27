# https://www.rubydoc.info/gems/shell/0.7/

Shell.alias_command("foo", params[:oops]) {}
Shell::CommandProcessor.alias_command(x, "/bin/#{params[:oops]}") {}

Shell.def_system_command("foo", "bar", params[:oops]) {}
Shell::CommandProcessor.def_system_command("foo", params[:oops]) {}

shell = Shell.new(Dir.pwd)

processor1 = shell.command_processor
processor1.system(params[:oops])

processor2 = Shell::CommandProcessor.new(shell)
processor2.system(x, params[:oops])
