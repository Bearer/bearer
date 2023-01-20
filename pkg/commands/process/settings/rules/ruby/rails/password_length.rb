# trigger:length is less than minimum limit
class User < ApplicationRecord
  has_secure_password
  validates :password, length: { minimum: 6 }
end

# ok:length is within minimum limit
class Student < ApplicationRecord
  has_secure_password
  validates :password, length: { minimum: 12 }
end