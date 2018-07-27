const serviceClient = require('./serviceClient')

module.exports = (emitter, event, { DeveloperPortalAPIUrl, orgId, scenarioId }) =>
  new Promise(async (resolve, reject) => {
    const client = serviceClient(DeveloperPortalAPIUrl)
    try {
      const res = await client.deployScenario(event, orgId, scenarioId)

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
