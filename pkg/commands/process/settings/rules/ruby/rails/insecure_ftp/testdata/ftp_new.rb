require "net/ftp"

# critical risk: application has sensitive data
class User
  attr_reader :name, :email, :password
end

ftp = Net::FTP.new("ftp.ruby-lang.org")
ftp.login("anonymous", "matz@ruby-lang.org")
ftp.chdir("/pub/ruby")
tgz = ftp.list("ruby-*.tar.gz").sort.last
ftp.getbinaryfile(tgz, tgz)
ftp.close