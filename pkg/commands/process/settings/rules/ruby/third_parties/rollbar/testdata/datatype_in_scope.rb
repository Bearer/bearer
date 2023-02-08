Rollbar.scope!({ user: { email: "someone@example.com" }})

user = { email: "someone@example.com" }

notifier = Rollbar.scope(user)

notifier.scope(user: { first_name: "someone" })
