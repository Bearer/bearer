import * as tsNode from 'ts-node'
import * as path from 'path'
import * as fs from 'fs'
import { parse } from 'jsonc-parser'

import { TAuthContext } from '@bearer/functions/lib/declaration'

import baseCommand from '../base-command'

const invokeFunction = async (command: baseCommand, funcName: string, body: any, extraContext: any) => {
  tsNode.register({
    project: path.join(command.locator.srcFunctionsDir, 'tsconfig.json')
  })
  let func
  try {
    func = require(path.join(command.locator.srcFunctionsDir, funcName)).default
  } catch (e) {
    throw e
  }

  return await func.init()({
    body,
    context: {
      ...getFunctionContext(command),
      ...extraContext
    }
  })
}

const getFunctionContext = (command: baseCommand) => {
  const localConfig = command.locator.localConfigPath
  const context = {} as TAuthContext
  if (fs.existsSync(localConfig)) {
    const rawConfig = fs.readFileSync(localConfig, { encoding: 'utf8' })
    const parsed = parse(rawConfig)
    const { setup } = parsed || { setup: null }
    // @ts-ignore
    command.debug('local config: %j', parsed)
    if (setup && setup.auth) {
      context.authAccess = setup.auth
    }
  } else {
    // @ts-ignore
    command.debug('no local config found')
    context.authAccess = {}
  }

  return context
}

export default invokeFunction
