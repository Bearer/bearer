# Insecure FTP

class User
  attr_reader :name, :email, :password
end

## Detected
require "net/ftp"

ftp = Net::FTP::new("ftp.ruby-lang.org")
ftp.login("anonymous", "matz@ruby-lang.org")
ftp.chdir("/pub/ruby")
tgz = ftp.list("ruby-*.tar.gz").sort.last
ftp.getbinaryfile(tgz, tgz)
ftp.close

## Not detected
require "net/sftp"

Net::SFTP.start("localhost", "user") do |sftp|
  sftp.upload! "/local/file.tgz", "/remote/file.tgz"
end
