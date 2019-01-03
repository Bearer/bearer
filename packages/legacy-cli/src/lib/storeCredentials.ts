const fs = require('fs')
const serviceClient = require('./serviceClient')

const readConfig = configFile =>
  new Promise<{ clientID: string; clientSecret: string }>((resolve, reject) =>
    fs.readFile(configFile, (err, data) => {
      if (err) reject(err)
      else resolve(JSON.parse(data))
    })
  )

export default async (configFile, { IntegrationServiceUrl }, emitter) => {
  try {
    const { clientID, clientSecret } = await readConfig(configFile)
    const client = serviceClient(IntegrationServiceUrl)
    if (clientID && clientSecret) {
      const {
        Item: { referenceId }
      } = await client.putItem({
        clientID,
        clientSecret
      })
      emitter.emit('storeCredentials:success', referenceId)
    } else {
      emitter.emit('storeCredentials:missingCredentials', configFile)
    }
  } catch (e) {
    emitter.emit('storeCredentials:failure', e)
  }
}
