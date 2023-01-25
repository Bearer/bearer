NewRelic::Agent.add_custom_attributes(other: 42)
NewRelic::Agent.add_custom_parameters(foo: "bar")
NewRelic::Agent.notice_error(exception)
