fd = IO.sysopen("/dev/tty", "w")
IO.open(fd,"w") do |a|
  a.puts "Hello, #{user.full_name}!"
end
