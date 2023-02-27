def my_handler(event:, context:)
  Dir["foo", base: event["oops"]]

  Dir.chdir("/home/#{event["oops"]}")

  File.exist?(event["oops"])

  IO.readlines("/home/#{event["oops"]}")

  Kernel.open(event["oops"], "w+") do
  end

  open(event["oops"])

  PStore.new(event["oops"])

  path = Pathname.new(event["oops"])
  path + event["two"]
  path / event["two"]
  path.join("a", event["three"])
end
