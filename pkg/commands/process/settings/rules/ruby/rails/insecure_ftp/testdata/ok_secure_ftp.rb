require "net/sftp"

Net::SFTP.start("localhost", "user") do |sftp|
  sftp.upload! "/local/file.tgz", "/remote/file.tgz"
end