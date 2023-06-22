matched_variable = 42
ignored = 42

class MatchedClass
end

class MatchedClass
  validates :password, length: { minimum: 2 }
end

class IgnoredClass
end

class IgnoredClass
  validates :password, length: { minimum: 2 }
end
