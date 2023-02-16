import * as opentelemetry from "@opentelemetry/sdk-node"

const tracer = opentelemetry.trace.getTracer("some-tracer");

tracer.startActiveSpan("app.some-span", { attributes: { name: currentUser.name } }, span => {
  if (somethingWentWrong()) {
    span.setStatus({
      code: opentelemetry.SpanStatusCode.ERROR,
      message: "Bad connection for " + currentUser.emailAddress
    })
  }

  span.end();
});