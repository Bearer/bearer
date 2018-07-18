const serviceClient = require('./serviceClient')
const fs = require('fs')

module.exports = (
  packagePath,
  { Key },
  emitter,
  {
    bearerConfig: {
      authorization: {
        AuthenticationResult: { IdToken: token }
      }
    },
    DeploymentUrl
  }
) =>
  new Promise(async (resolve, reject) => {
    emitter.emit('pushScenario:start', Key)

    try {
      const deploymentServiceClient = serviceClient(DeploymentUrl)

      const res = await deploymentServiceClient.signedUrl(token, Key, 'intent')

      if (res.statusCode === 201) {
        const url = res.body
        const s3Client = serviceClient(url)
        const artifact = fs.readFileSync(packagePath)
        const response = await s3Client.upload(artifact)
        resolve(response)
      } else if (res.statusCode === 401) {
        emitter.emit('pushScenario:unauthorized', res.body)
        reject(new Error('unauthorized'))
      } else {
        emitter.emit('pushScenario:httpError', res)
        reject(new Error('httpError'))
      }
    } catch (e) {
      emitter.emit('pushScenario:error', e)
      reject(e)
    }
  })
