const fs = require('fs') // from node.js
const path = require('path') // from node.js
const globby = require('globby')
const mime = require('mime-types')
const serviceClient = require('./serviceClient')

const DIST_DIRECTORY = 'dist'
const WWW_DIRECTORY = 'www'

const pushViews = (
  viewsDirectory,
  emitter,
  {
    DeploymentUrl,
    scenarioId,
    orgId,
    DeveloperPortalAPIUrl,
    bearerConfig: {
      authorization: {
        AuthenticationResult: { IdToken: token }
      }
    },
    credentials
  }
) => {
  return new Promise(async (resolve, reject) => {
    const configuration = {
      distPath: path.join(viewsDirectory, DIST_DIRECTORY),
      wwwPath: path.join(viewsDirectory, WWW_DIRECTORY)
    }

    const integrationsClient = serviceClient(DeploymentUrl)
    const devPortalClient = serviceClient(DeveloperPortalAPIUrl)
    const {
      body: {
        data: {
          findUser: { token: devPortalToken }
        }
      }
    } = await devPortalClient.getDevPortalToken(credentials)
    try {
      emitter.emit('view:upload:start')

      const files = await globby([configuration.distPath, configuration.wwwPath])

      const paths = files.reduce((acc, filePath) => {
        const relativePath = filePath.replace(viewsDirectory + path.sep, '')
        acc[`${orgId}/${scenarioId}/${relativePath}`] = filePath
        return acc
      }, {})

      const urls = (await integrationsClient.signedUrls(token, Object.keys(paths), 'screen')).body
      await Promise.all(
        Object.keys(paths).map(key =>
          (async () => {
            try {
              const filePath = paths[key]
              const fileContent = fs.readFileSync(filePath)
              const s3Client = serviceClient(urls[key])
              await s3Client.upload(fileContent.toString(), {
                'Content-Type': mime.lookup(filePath)
              })
            } catch (e) {
              emitter.emit('view:fileUpload:error', e)
              reject(e)
            }
          })()
        )
      )
      resolve('done')
    } catch (e) {
      emitter.emit('view:upload:error', e)
      reject(e)
    }
  })
}

module.exports = pushViews
