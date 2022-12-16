# it should match since max is lower than required
class User < ApplicationRecord
	devise password_length: 11..32
end

# it should match since minimum is lower than required
class Employee < ApplicationRecord
	validates :password, length: { minimum: 6 }
end

# it should match since min is lower than required
class Student < ApplicationRecord
	devise password_length: 6..36
end

# it shouldn't match since max and min are withing boundaries
class Admin < ApplicationRecord
	devise password_length: 11..36
end