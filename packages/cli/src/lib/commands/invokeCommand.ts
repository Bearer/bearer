import axios from 'axios'
import * as querystring from 'querystring'
import * as cosmiconfig from 'cosmiconfig'
import * as jq from 'node-jq'

const startLocalDevelopmentServer = require('./startLocalDevelopmentServer')

import Locator from '../locationProvider'

type IConfig = {
  httpMethod?: string
  params?: {}
  body?: {}
}

export const invoke = (emitter, config, locator: Locator) => async (intent, cmd) => {
  const { path, format } = cmd
  const {
    scenarioUuid,
    scenarioConfig: { scenarioTitle }
  } = config
  const DEFAULT_JQ_FORMAT = '.'
  const jqFilter = format ? format : DEFAULT_JQ_FORMAT
  const jqOptions = { input: 'string' }

  let fileData: IConfig = {}
  if (path) {
    const explorer = cosmiconfig(path, {
      searchPlaces: [path]
    })
    const { config = {} } = (await explorer.search(locator.scenarioRootResourcePath.toString())) || {}
    fileData = config
  }
  const { httpMethod = 'GET', params = {}, body = {} } = fileData

  const integrationHostURL = await startLocalDevelopmentServer(scenarioUuid, emitter, config, locator)

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
    jq.run(jqFilter, JSON.stringify(intentData), jqOptions)
      .then(output => {
        console.log(output)
        process.exit(0)
      })
      .catch(err => {
        console.log(intentData)
        console.log(err)
        process.exit(1)
      })
  } catch (e) {
    console.log(e)
    process.exit(1)
  }
}

export function useWith(program, emitter, config, locator: Locator) {
  program
    .command('invoke <intent>')
    .option('-p, --path <path>')
    .option('-f, --format <format>')
    .description(
      `invoke Intent locally.
  $ bearer invoke <IntentName>
  $ bearer invoke <IntentName> -p tests/intent.json
  $ bearer invoke <IntentName> -f '.data[].name' ## JQ compatible format
`
    )
    .action(invoke(emitter, config, locator))
}
