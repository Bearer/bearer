const inquirer = require('inquirer')
const serviceClient = require('../serviceClient')

const login = (emitter, config) => async ({ email }) => {
  let {
    bearerConfig: { Username }
  } = config
  if (!Username && !email) {
    emitter.emit('username:missing')
    process.exit(1)
  }

  if (email) Username = email

  emitter.emit('login:userFound', Username)
  try {
    const client = serviceClient(config.IntegrationServiceUrl)
    const { AccessToken } = await inquirer.prompt([
      {
        message: `Please enter your access token:`,
        type: 'password',
        name: 'AccessToken'
      }
    ])

    const res = await client.login({ Username, Password: AccessToken })

    let ExpiresAt
    switch (res.statusCode) {
      case 200:
        ExpiresAt = res.body.authorization.AuthenticationResult.ExpiresIn + Date.now()
        config.storeBearerConfig({
          ...res.body.user,
          ExpiresAt,
          authorization: res.body.authorization
        })
        emitter.emit('login:success', res.body)
        break
      case 401:
        emitter.emit('login:failure', res.body)
        break
      default:
        emitter.emit('login:error', { code: res.statusCode, body: res.body })
    }
  } catch (e) {
    emitter.emit('login:error', e)
  }
}
module.exports = {
  useWith: (program, emitter, config) => {
    program
      .command('login')
      .description(
        `Login to Bearer.
    $ bearer login
`
      )
      .option('-e, --email <email>', 'User email.')
      .action(login(emitter, config))
  }
}
