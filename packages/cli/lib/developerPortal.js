const serviceClient = require('./serviceClient')

module.exports = (
  emitter,
  event,
  {
    DeveloperPortalAPIUrl,
    bearerConfig: { OrgId },
    scenarioConfig: { scenarioTitle }
  }
) =>
  new Promise(async (resolve, reject) => {
    const client = serviceClient(DeveloperPortalAPIUrl)
    try {
      const res = await client.deployScenario(event, OrgId, scenarioTitle)

      if (!res.errors) {
        emitter.emit('developerPortalUpdate:success')
      } else {
        emitter.emit('developerPortalUpdate:failed', res.errors)
      }

      resolve('done')
    } catch (e) {
      emitter.emit('developerPortalUpdate:error', e)
      reject(e)
    }
  })
