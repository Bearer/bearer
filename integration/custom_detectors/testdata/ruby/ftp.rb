# No data type detection expected
Net::FTP.open("ftp.site.com") do |ftp|
  ftp.puttextfile("no_data.txt", "/no_data.txt")
end


user = { email: "dave@example.com" }

# Expecting data type detection (but not working yet)
Net::FTP.open("ftp.site.com") do |ftp|
  file = Tempfile.new("user_data")
  begin
    file << user.to_json
    file.close

    ftp.puttextfile(file.path, "/users/123.json")
  ensure
    file.close!
  end
end


# Expecting data type detection
Net::FTP.open("ftp.site.com") do |ftp|
  file = Tempfile.new("user_data")
  begin
    file << { user: { ethnicity: "martian" } }.to_json
    file.close

    ftp.puttextfile(file.path, "/users/123.json")
  ensure
    file.close!
  end
end
