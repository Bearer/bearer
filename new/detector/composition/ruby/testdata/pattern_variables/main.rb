matched_variable = 42
ignored = 42

class MatchedClass
end

# the error parsing is different between classes with bodies and without
class MatchedClass
  something :foo
end

class IgnoredClass
end
