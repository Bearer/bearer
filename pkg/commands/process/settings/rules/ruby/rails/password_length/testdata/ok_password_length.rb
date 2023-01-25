class Student < ApplicationRecord
  has_secure_password
  validates :password, length: { minimum: 12 }
end