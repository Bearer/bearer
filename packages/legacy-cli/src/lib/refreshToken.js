const serviceClient = require('./serviceClient')

module.exports = async (config, emitter) => {
  const client = serviceClient(config.IntegrationServiceUrl)
  const { RefreshToken } = config.bearerConfig.authorization.AuthenticationResult
  let { bearerConfig } = config
  try {
    const res = await client.refresh({ RefreshToken })

    let ExpiresAt
    let AccessToken
    let IdToken
    let ExpiresIn

    switch (res.statusCode) {
      case 200:
        ExpiresAt = res.body.authorization.AuthenticationResult.ExpiresIn + Date.now()

        AccessToken = res.body.authorization.AuthenticationResult.AccessToken
        IdToken = res.body.authorization.AuthenticationResult.IdToken
        ExpiresIn = res.body.authorization.AuthenticationResult.ExpiresIn

        bearerConfig = Object.assign(config.bearerConfig, {
          ExpiresAt,
          authorization: {
            AuthenticationResult: {
              AccessToken,
              IdToken,
              ExpiresIn,
              RefreshToken
            }
          }
        })
        config.storeBearerConfig(bearerConfig)
        emitter.emit('refreshToken:success', res.body)
        break
      case 401:
        emitter.emit('refreshToken:failure', res.body)
        break
      default:
        emitter.emit('refreshToken:error', {
          code: res.statusCode,
          body: res.body
        })
    }
  } catch (e) {
    emitter.emit('refreshToken:error', e)
  }

  return Object.assign(config, { bearerConfig })
}
