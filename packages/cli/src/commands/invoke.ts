import { flags } from '@oclif/command'

import * as path from 'path'
import * as tsNode from 'ts-node'
import * as cosmiconfig from 'cosmiconfig'
import BaseCommand from '../base-command'

const LOCAL_DEV_CONFIGURATION = 'dev'

export default class Invoke extends BaseCommand {
  static description = 'Invoke Function locally'

  static flags = {
    ...BaseCommand.flags,
    data: flags.string({ char: 'd' }),
    file: flags.string({ char: 'f' })
    // TODO: add arguments : allow file or daja
  }

  static args = [{ name: 'Function_Name', required: true }]

  async run() {
    const { args } = this.parse(Invoke)
    const funcName = args.Function_Name
    this.debug(funcName)
    tsNode.register({
      project: path.join(this.locator.srcFunctionsDir, 'tsconfig.json')
    })
    let func
    try {
      func = require(path.join(this.locator.srcFunctionsDir, funcName)).default
    } catch (e) {
      const funcName = args.Function_Name
      if (e.code === 'MODULE_NOT_FOUND') {
        // TODO: add information, check file presence or generator
        this.error(`"${funcName}" function does not exist. `)
      }
      this.error(e.message)
    }

    const explorer = cosmiconfig(LOCAL_DEV_CONFIGURATION, {
      searchPlaces: [`config.${LOCAL_DEV_CONFIGURATION}.js`]
    })
    const { config: devFunctionsContext = {} } = (await explorer.search(this.locator.integrationRoot)) || {}
    const bearerBaseURL = 'ok'
    // TODO: access sqlite database and load defined data
    const userDefinedData = {}
    const datum = await func.init()(
      {
        context: {
          ...devFunctionsContext.global,
          ...devFunctionsContext[funcName],
          bearerBaseURL,
          ...userDefinedData
        },
        body: JSON.stringify({})
      },
      {}
    )

    console.log(JSON.stringify(datum, null, 2))
  }
}
