const { Notifier } = require('@airbrake/node');

const airbrake = new Notifier({
  projectId: 42,
  projectKey: 'some-project-key',
  environment: 'PROD',
});

airbrake.notify({
  error: err,
  params: { env: "prod" },
});