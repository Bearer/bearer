const serviceClient = require('./serviceClient')

module.exports = (
  emitter,
  {
    DeploymentUrl,
    bearerConfig: { OrgId },
    scenarioConfig: { scenarioTitle },
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
      const invalidationPath = `${OrgId}/${scenarioTitle}`
      const res = await client.viewsInvalidate(token, { invalidationPath })

      if (res.statusCode === 204) emitter.emit('invalidateCloudFront:success')

      if (res.statusCode !== 204) emitter.emit('invalidateCloudFront:invalidationFailed', res.body)

      resolve('done')
    } catch (e) {
      emitter.emit('invalidateCloudFront:error', e)
      reject(e)
    }
  })
