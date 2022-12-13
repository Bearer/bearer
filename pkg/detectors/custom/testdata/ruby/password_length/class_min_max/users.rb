# it should match since max is lower than required
class User < ApplicationRecord
	device password_length: 11..32
end


# it should match since min is lower than required
class Student < ApplicationRecord
	device password_length: 6..36
end

# it shouldn't match since max and min are withing boundaries
class Admin < ApplicationRecord
	device password_length: 11..36
end