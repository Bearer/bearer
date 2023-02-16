Net::FTP.new(params[:oops])

Net::FTP.open("example.com", username: params[:user]) do

end

def handler(event:, context:)
  ftp = Net::FTP.open("example.com")
  ftp.puttextfile("local.txt", event["filename"])
end
