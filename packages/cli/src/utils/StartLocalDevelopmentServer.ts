import * as express from 'express'
import * as getPort from 'get-port'
import * as morgan from 'morgan'

import Start from '../commands/start'

import LocationProvider from './locator'

const app: express.Express = express()

app.use(morgan('dev'))

const port = 3000

const intentRouter: (locator: LocationProvider) => express.Router = (locator: LocationProvider) => {
  const app = express.Router()

  function requireUncached(m: string) {
    delete require.cache[require.resolve(m)]
    return require(m)
  }

  app.all('/:scenarioUuid/:intentName', (req: express.Request, res: express.Response) => {
    const { intentName } = req.params
    const { default: { intentType, action } } = requireUncached(`${locator.buildIntentsDir}/dist/${intentName}`)

    intentType.intent(action)({
      context: {},
      queryStringParameters: req.query,
      body: req.body
    },
      {},
      (_err: any, datum: any) => {
        res.json(datum)
      })

  })

  return app
}


export default {
  run: async (cmd: Start) => {
    const openPort = await getPort({ port })
    app.use(intentRouter(cmd.locator))
    app.listen(openPort)
  }
}
