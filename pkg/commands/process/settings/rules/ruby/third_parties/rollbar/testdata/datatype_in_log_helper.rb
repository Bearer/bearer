Rollbar.critical("oops #{user.email}")
Rollbar.critical(e, "oops #{user.email}")
Rollbar.critical(e, user: { email: "someone@example.com" })
Rollbar.critical(e, { user: { first_name: "someone" } })

Rollbar.error("oops #{user.email}")

Rollbar.debug("oops #{user.email}")

Rollbar.info("oops #{user.email}")

Rollbar.warning("oops #{user.email}")
