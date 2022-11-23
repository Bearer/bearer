class Users < ApplicationRecord
    encrypts :email, :country, :city # requires the detection of the structure too
end