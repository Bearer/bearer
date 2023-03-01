if request.get?
elsif request.post?
  change_state
end

change_state if request.post?

if request.post?
  change_state
end
