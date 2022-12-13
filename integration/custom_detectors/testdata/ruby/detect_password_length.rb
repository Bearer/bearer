# it should match this since minimum is below needed one
class User < ApplicationRecord
    has_secure_password
    validates :password, length: { minimum: 6 }
end

# it shouldn't match this since minimum is above needed one
class Student < ApplicationRecord
    has_secure_password
    validates :password, length: { minimum: 12 }
end