exception.rollbar_context = { foo: "bar" }

Rollbar.critical("oops")

Rollbar.log("error", "oops")

Rollbar.scope!( { foo: "bar" })

scope = { foo: "bar" }

Rollbar.scoped(scope) do
  call
end
