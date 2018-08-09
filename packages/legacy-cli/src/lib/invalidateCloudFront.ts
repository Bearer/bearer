import * as serviceClient from './serviceClient'
import { Config } from './types'

export default (
  emitter,
  {
    DeploymentUrl,
    scenarioConfig: { orgId, scenarioId },
    bearerConfig: {
      authorization: {
        AuthenticationResult: { IdToken: token }
      }
    }
  }: Config
) =>
  new Promise(async (resolve, reject) => {
    const client = serviceClient(DeploymentUrl)
    try {
      const invalidationPath = `${orgId}/${scenarioId}`
      const res = await client.viewsInvalidate(token, { invalidationPath })

      if (res.statusCode === 204) emitter.emit('invalidateCloudFront:success')

      if (res.statusCode !== 204) emitter.emit('invalidateCloudFront:invalidationFailed', res.body)

      resolve('done')
    } catch (e) {
      emitter.emit('invalidateCloudFront:error', e)
      reject(e)
    }
  })
