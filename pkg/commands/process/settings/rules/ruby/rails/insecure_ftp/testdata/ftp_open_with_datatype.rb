require "net/ftp"

Net::FTP.open("ftp.site.com") do |ftp|
  file = Tempfile.new("user_data")
  begin
    file << user.email
    file.close

    ftp.puttextfile(file.path, "/users/123.json")
  ensure
    file.close!
  end
end