# frozen_string_literal: true
class User < Db1Record
    include Corruptible
  
  
    def corruptible_attributes_by_type
      {
        CorruptibleType::STRING => %i[
          email
          address
        ],
        CorruptibleType::ID => %i[first_name last_name middle_name],
        CorruptibleType::DATE => %i[birthdate],
      }
    end
end