Rollbar.log("error", "oops #{user.email}")
Rollbar.log("error", "oops", user: { email: "someone@example.com" })
Rollbar.log("error", "oops", { user: { first_name: "someone" } })
