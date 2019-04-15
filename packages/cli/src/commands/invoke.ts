import { flags } from '@oclif/command'

import * as path from 'path'
import * as fs from 'fs'
import * as tsNode from 'ts-node'
import * as http from 'http'

import { parse } from 'jsonc-parser'
import { TAuthContext } from '@bearer/functions/lib/declaration'
import { TPersistedData } from '@bearer/functions/lib/db-client'

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

    const body = flags.data || (flags.file && this.getFileContent(flags.file)) || '{}'
    this.debug('Injected body, %j', body)

    const context = this.getFunctionContext()
    // injected within the function (let's remove it?)
    const bearerBaseURL = 'ok'

    this.ensureJson(body)

    this.debug('calling')
    process.env.bearerBaseURL = 'http://localhost:5678'
    this.debug('starting')
    const a: any = {}
    const server: http.Server = await new Promise((resolve, reject) => {
      const _server = http
        .createServer((request, response) => {
          this.debug('Incoming request')
          let body = ''
          request.on('data', chunk => {
            body += chunk
          })
          request.on('end', () => {
            request.method
            try {
              // const data: { code: string } = JSON.parse(body)
              this.debug('method: %s body: %s', request.method, body || '{}')
              response.setHeader('Connection', 'close')

              switch (request.method) {
                case 'GET': {
                  const referenceId = request.url!.split('/').reverse()[0]
                  const existingData = a[referenceId]
                  this.debug('lookup id: %s data:%j', referenceId, existingData)
                  const payload: TPersistedData = { Item: { referenceId, ReadAllowed: true } }
                  if (existingData) {
                    payload.Item.data = existingData
                  }
                  this.debug('GET sending %j', payload)
                  response.write(JSON.stringify(payload))
                  break
                }
                case 'POST':
                case 'PUT': {
                  const { referenceId, data } = JSON.parse(body)
                  const returnedData: TPersistedData = { Item: { ...data } }
                  a[referenceId] = data
                  this.debug('perisisting data %j', data)
                  response.write(JSON.stringify(returnedData))
                }
              }
            } catch (e) {
              this.debug('error:', e)
            }
            response.end()
          })
        })
        .listen(5678, () => {
          this.debug('started')
          resolve(_server)
        })
    })

    const datum = await func.init()({
      body,
      context: {
        ...context,
        bearerBaseURL,
        isBackend: true // TODO: flag option
      }
    })
    server.close()
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

  getFileContent = (filePath: string): string | undefined => {
    const location = path.resolve(filePath)
    if (filePath) {
      return fs.existsSync(location)
        ? fs.readFileSync(location, { encoding: 'utf8' })
        : this.error(`File not found: ${location}`)
    }
    return undefined
  }

  getFunctionContext() {
    const localConfig = this.locator.localConfigPath
    const context = {} as TAuthContext
    if (fs.existsSync(localConfig)) {
      const rawConfig = fs.readFileSync(localConfig, { encoding: 'utf8' })
      const parsed = parse(rawConfig)
      const { setup } = parsed || { setup: null }
      this.debug('local config: %j', parsed)
      if (setup && setup.auth) {
        context.authAccess = setup.auth
      }
    } else {
      this.debug('no local config found')
      context.authAccess = {}
    }

    return context
  }
}
