Net::FTP.new(x)

Net::FTP.open("example.com", username: "user") do

end

event = not_from_handler
ftp = Net::FTP.open("example.com")
ftp.puttextfile("local.txt", event["filename"])
