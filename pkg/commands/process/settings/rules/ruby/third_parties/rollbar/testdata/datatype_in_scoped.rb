scope = { person: { email: "someone@example.com" } }

Rollbar.scoped(scope) do
  call
end
