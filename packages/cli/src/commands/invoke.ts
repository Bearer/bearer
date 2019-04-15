import { flags } from '@oclif/command'

import * as path from 'path'
import * as fs from 'fs'
import * as tsNode from 'ts-node'
import BaseCommand from '../base-command'

export default class Invoke extends BaseCommand {
  static description = 'invoke function locally'

  static flags = {
    ...BaseCommand.flags,
    data: flags.string({ char: 'd' }),
    file: flags.string({ char: 'f' })
    // TODO: add arguments : allow file or daja
  }

  static args = [{ name: 'Function_Name', required: true }]

  async run() {
    const { args, flags } = this.parse(Invoke)
    const funcName = args.Function_Name
    this.debug('start invoking: %j', funcName)
    tsNode.register({
      project: path.join(this.locator.srcFunctionsDir, 'tsconfig.json')
    })
    this.debug('ts-node registered')
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

    this.debug('required')
    // inject auth data
    // use body
    // inject required data

    this.debug('injecting data')
    const body = flags.data || (flags.file && this.getFilecontent(flags.file)) || '{}'
    const config = {} as any
    // injected within the function (let's remove it?)
    const bearerBaseURL = 'ok'
    // TODO: access sqlite database and load defined data from params sent
    const userDefinedData = {}
    this.ensureJson(body)

    this.debug('calling')
    const datum = await func.init()({
      body,
      context: {
        ...config.global,
        ...config[funcName],
        bearerBaseURL,
        ...userDefinedData
      }
    })
    // print output
    console.log(JSON.stringify(datum, null, 2))
  }

  ensureJson = (maybeJson: string) => {
    try {
      JSON.parse(maybeJson)
    } catch (e) {
      this.error('Invalid JSON provided')
    }
  }

  getFilecontent = (filePath: string): string | undefined => {
    const location = path.resolve(filePath)
    if (filePath) {
      return fs.existsSync(location)
        ? fs.readFileSync(location, { encoding: 'utf8' })
        : this.error(`File not found: ${location}`)
    }
    return undefined
  }
}
