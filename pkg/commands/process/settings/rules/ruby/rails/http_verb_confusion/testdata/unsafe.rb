if request.get?
else
  change_state
end

change_state unless request.get?

unless request.get?
  change_state
end
