import * as opentelemetry from "@opentelemetry/sdk-node"

var currentSpan = opentelemetry.trace.getSpan()
currentSpan.addEvent('my-event', {
  'event.metadata': customer.emailAddress
})
