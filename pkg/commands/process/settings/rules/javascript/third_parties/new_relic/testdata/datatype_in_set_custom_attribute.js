const newrelic = require("newrelic")

newrelic.setCustomAttribute("user-id", customer.email)
