import { flags } from '@oclif/command'

import * as path from 'path'
import getPort from 'get-port'
import * as fs from 'fs'
import * as tsNode from 'ts-node'
import * as http from 'http'
import * as knex from 'knex'

import { parse } from 'jsonc-parser'
import { TAuthContext } from '@bearer/functions/lib/declaration'
import { TPersistedData } from '@bearer/functions/lib/db-client'

import BaseCommand from '../base-command'
const filename = process.env.BEARER_LOCAL_DATABASE || ':memory:'

export default class Invoke extends BaseCommand {
  static description = 'invoke function locally'

  static flags = {
    ...BaseCommand.flags,
    data: flags.string({ char: 'd' }),
    file: flags.string({ char: 'f' })
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
        this.error(`"${funcName}" function does not exist. `)
      }
      this.error(e.message)
    }

    this.debug('required')

    const body = flags.data || (flags.file && this.getFileContent(flags.file)) || '{}'
    this.debug('Injected body, %j', body)

    const context = this.getFunctionContext()

    this.ensureJson(body)

    this.debug('calling')
    const port = await getPort()
    const bearerBaseURL = `http://localhost:${port}`
    process.env.bearerBaseURL = bearerBaseURL
    this.debug('starting')
    const server: http.Server = await this._startServer(port)

    const datum = await func.init()({
      body,
      context: {
        ...context,
        bearerBaseURL,
        isBackend: true // TODO: flag option
      }
    })
    server.close()
    if (this._db) {
      this._db.destroy()
    }

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

  _startServer = async (port: number) => {
    return new Promise<http.Server>((resolve, reject) => {
      const _server = http
        .createServer((request, response) => {
          this.debug('Incoming request')
          let body = ''
          request.on('data', chunk => {
            body += chunk
          })
          request.on('end', async () => {
            request.method
            try {
              this.debug('method: %s body: %s', request.method, body || '{}')
              response.setHeader('Connection', 'close')

              switch (request.method) {
                case 'GET': {
                  const referenceId = request.url!.split('/').reverse()[0]
                  const existingData = await this.getData(referenceId)
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
                  await this.putData(referenceId, data)
                  this.debug('persisting data %j', data)
                  const returnedData: TPersistedData = { Item: { ...data } }
                  response.write(JSON.stringify(returnedData))
                }
              }
            } catch (e) {
              this.debug('error:', e)
            }
            response.end()
          })
        })
        .listen(port, () => {
          this.debug('started')
          resolve(_server)
        })
    })
  }

  private _db!: knex

  async db(): Promise<knex.QueryInterface> {
    if (!this._db) {
      this._db = knex({
        dialect: 'sqlite3',
        connection: {
          filename
        },
        pool: { min: 0 },
        useNullAsDefault: true,
        debug: /bearer:/.test(process.env.DEBUG || '')
      })
      await this._db.schema
        .hasTable('records')
        .then(async (exists: boolean) => {
          if (!exists) {
            await this._db.schema.createTable('records', (table: any) => {
              table.index(['referenceId'], 'ref_id_idx')
              table.string('referenceId')
              table.text('data')
            })
          }
        })
        .catch(() =>
          this._db.schema.createTable('records', (table: any) => {
            table.index(['referenceId'], 'ref_id_idx')
            table.string('referenceId')
            table.text('data')
          })
        )
    }
    return this._db
  }

  getData = async (referenceId: string) => {
    const db = await this.db()
    const rows = await db
      .table('records')
      .select('data')
      .where({ referenceId })
      .limit(1)

    if (rows[0] as any) {
      return JSON.parse(rows[0].data)
    }
    return {}
  }

  putData = async (referenceId: string, data: any) => {
    const db = await this.db()
    db.table('records')
      .where({ referenceId })
      .delete()
    await db.table('records').insert({ referenceId, data: JSON.stringify(data) })
  }
}
