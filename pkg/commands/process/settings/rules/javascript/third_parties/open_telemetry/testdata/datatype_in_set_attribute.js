import * as opentelemetry from "@opentelemetry/sdk-node"

const tracer = opentelemetry.trace.getTracer("some-tracer");

tracer.startActiveSpan("some-span", span => {
  span.setAttribute("current-user", currentUser.emailAddress);
  span.end();
});


span.setAttribute("current-user", user.email);
