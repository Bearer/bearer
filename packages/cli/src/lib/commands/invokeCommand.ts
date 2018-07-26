import axios from 'axios'
import * as querystring from 'querystring'
import * as cosmiconfig from 'cosmiconfig'

const startLocalDevelopmentServer = require('./startLocalDevelopmentServer')

import Locator from '../locationProvider'

class NoEmitter {
  emit(_name, _args) {}

  on(_args) {}
}

export const invoke = (_emitter, config, locator: Locator) => async (intent, cmd) => {
  const { file } = cmd
  const explorer = cosmiconfig(file, {
    searchPlaces: [file]
  })
  const {
    bearerConfig: { OrgId },
    scenarioConfig: { scenarioTitle }
  } = config

  const { config: fileData = {} } = (await explorer.search(locator.scenarioRootResourcePath.toString())) || {}
  const { httpMethod = 'GET', params = {}, body = {} } = fileData

  const scenarioUuid = `${OrgId}-${scenarioTitle}`
  const noEmitter = new NoEmitter()
  const integrationHostURL = await startLocalDevelopmentServer(scenarioUuid, noEmitter, config, locator, false)

  const client = axios.create({
    baseURL: `${integrationHostURL}api/v1`,
    timeout: 5000
  })

  try {
    let intentData
    if (httpMethod == 'GET') {
      const { data } = await client.get(`${scenarioUuid}/${intent}?${querystring.stringify(params)}`)
      intentData = data
    }
    if (httpMethod == 'POST') {
      const { data } = await client.post(`${scenarioUuid}/${intent}`, querystring.stringify(body))
      intentData = data
    }
    console.log(JSON.stringify(intentData))
    process.exit(0)
  } catch (e) {
    console.log(e)
    process.exit(0)
  }
}

export function useWith(program, emitter, config, locator: Locator) {
  program
    .command('invoke <intent>')
    .option('-f, --file <path>')
    .description(
      `invoke Intent locally.
  $ bearer invoke <IntentName>
`
    )
    .action(invoke(emitter, config, locator))
}
