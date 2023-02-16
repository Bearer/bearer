import * as opentelemetry from "@opentelemetry/sdk-node"

var span = opentelemetry.trace.getSpan()

try {
  // something
} catch (err) {
  span.recordException(err)
  span.recordException(currentUser.ipAddress)
}