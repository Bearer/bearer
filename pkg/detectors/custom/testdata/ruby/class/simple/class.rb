# Detected
class User < ApplicationRecord
    encrypts :email
end

# Not detected
class User < ApplicationRecord
    something :email
end

class User
    something :email
end

class User
    encrypts :email
end

class User < Base
    encrypts :email
end

class User < Base
    attr_reader :email
    attr_writer :email
    attr_accessor :email
end