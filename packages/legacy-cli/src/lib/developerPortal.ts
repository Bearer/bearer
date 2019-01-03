import serviceClient from './serviceClient'

export default (
  emitter,
  event,
  {
    DeploymentUrl,
    orgId,
    scenarioId,
    bearerConfig: {
      authorization: {
        AuthenticationResult: { IdToken: token }
      }
    }
  }
) =>
  new Promise(async (resolve, reject) => {
    const client = serviceClient(DeploymentUrl)
    try {
      const res = await client.deployScenario(token, event, orgId, scenarioId)
      if (res.statusCode === 204 || res.statusCode === 202 || res.statusCode === 200) {
        emitter.emit('developerPortalUpdate:success')
      } else {
        emitter.emit('developerPortalUpdate:failed', res.body.errors)
      }

      resolve('done')
    } catch (e) {
      emitter.emit('developerPortalUpdate:error', e)
      reject(e)
    }
  })
