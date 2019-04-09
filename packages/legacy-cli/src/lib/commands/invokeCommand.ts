import axios from 'axios'
import * as cosmiconfig from 'cosmiconfig'

import Locator from '../locationProvider'
import startLocalDevelopmentServer from './startLocalDevelopmentServer'

type IConfig = {
  params?: {}
  body?: {}
}

export const invoke = (emitter, config, locator: Locator) => async (func, cmd) => {
  const { path } = cmd
  const { buid } = config

  let fileData: IConfig = {}
  if (path) {
    const explorer = cosmiconfig(path, {
      searchPlaces: [path]
    })
    const { config = {} } = (await explorer.search(locator.integrationRootResourcePath.toString())) || {}
    fileData = config
  }
  const { params = {}, body = {} } = fileData
  try {
    const integrationHostURL = await startLocalDevelopmentServer(emitter, config, locator, { force: true })

    const client = axios.create({ baseURL: `${integrationHostURL}api/v4/backend/functions`, timeout: 5000 })

    const { data } = await client.post(`${buid}/${func}`, body, { params })
    // used by cli: do not remove
    console.log(JSON.stringify(data, null, 2))
    process.exit(0)
  } catch (e) {
    if (e.response) {
      console.log(JSON.stringify(e.response.data, null, 2))
      process.exit(0)
    } else {
      console.log(e.toString())
      process.exit(1)
    }
  }
}

export function useWith(program, emitter, config, locator: Locator) {
  program
    .command('invoke <func>')
    .option('-p, --path <path>')
    .description(
      `invoke Function locally.
  $ bearer invoke <FunctionName>
  $ bearer invoke <FunctionName> -p tests/func.json
`
    )
    .action(invoke(emitter, config, locator))
}
