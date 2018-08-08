import axios from 'axios'
import * as cosmiconfig from 'cosmiconfig'
import startLocalDevelopmentServer from './startLocalDevelopmentServer'

import Locator from '../locationProvider'

type IConfig = {
  params?: {}
  body?: {}
}

export const invoke = (emitter, config, locator: Locator) => async (intent, cmd) => {
  const { path } = cmd
  const { scenarioUuid } = config

  let fileData: IConfig = {}
  if (path) {
    const explorer = cosmiconfig(path, {
      searchPlaces: [path]
    })
    const { config = {} } = (await explorer.search(locator.scenarioRootResourcePath.toString())) || {}
    fileData = config
  }
  const { params = {}, body = {} } = fileData
  try {
    const integrationHostURL = await startLocalDevelopmentServer(emitter, config, locator)

    const client = axios.create({ baseURL: `${integrationHostURL}api/v1`, timeout: 5000 })

    const { data } = await client.post(`${scenarioUuid}/${intent}`, body, { params })
    console.log(JSON.stringify(data, null, 2))
    process.exit(0)
  } catch (e) {
    console.log(e)
    process.exit(1)
  }
}

export function useWith(program, emitter, config, locator: Locator) {
  program
    .command('invoke <intent>')
    .option('-p, --path <path>')
    .description(
      `invoke Intent locally.
  $ bearer invoke <IntentName>
  $ bearer invoke <IntentName> -p tests/intent.json
`
    )
    .action(invoke(emitter, config, locator))
}
