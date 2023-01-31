user = { first_name: "foo" }
NewRelic::Agent.notice_error(exception, { custom_params: user })
NewRelic::Agent.notice_error(exception, expected: true, custom_params: { last_name: "foo" }, metric: "test")
