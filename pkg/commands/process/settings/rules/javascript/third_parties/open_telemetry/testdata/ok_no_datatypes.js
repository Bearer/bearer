import * as opentelemetry from "@opentelemetry/sdk-node"

const tracer = opentelemetry.trace.getTracer("some-tracer");

tracer.startActiveSpan("some-span", span => {
  span.setAttribute("some-attribute", "some-value");
  span.end();
});

