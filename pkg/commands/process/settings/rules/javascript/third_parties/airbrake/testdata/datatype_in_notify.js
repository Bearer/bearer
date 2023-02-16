const { Notifier } = require('@airbrake/node');

const airbrake = new Notifier({
  projectId: 42,
  projectKey: 'some-project-key',
  environment: 'PROD',
});

let promise = airbrake.notify("user " + currentUser().emailAddress)
promise.then(() => {})

riskyCode(() => {
  try {
    // something risky
  } catch (err) {
    airbrake.notify({
      error: err,
      params: { user: user.ipAddress },
    });
  }
})