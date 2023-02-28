# https://www.rubydoc.info/gems/shell/0.7/

Shell.cd(params[:oops])

Shell.default_system_path = params[:oops]

shell = Shell.new(params[:oops], umask)

shell.pushdir(params[:oops], true)

processor1 = shell.command_processor
processor1.foreach(params[:oops], rs) {}

processor2 = Shell::CommandProcessor.new(shell)
processor2.test(:exists?, x, params[:oops])

processor2[:exists?, x, params[:oops], y]

processor2.transact do
  test(:exists?, params[:oops])
end

# optional arg
shell = Shell.new(params[:oops])