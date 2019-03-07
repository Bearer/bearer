import { Authentications } from '@bearer/types/lib/authentications'
import FunctionType from '@bearer/types/lib/function-types'
import * as Case from 'case'
import * as path from 'path'

import { copyFiles } from '../helpers'
import baseCommand from '../../base-command'

export default async (command: baseCommand, auth: Authentications, type: FunctionType, name: string) => {
  const vars = getVars(name, auth, type)
  return await copyFiles(command, path.join(`generate/intent`, auth, type), command.locator.srcFunctionsDir, vars)
}

function getVars(name: string, authType: Authentications, intentType: FunctionType) {
  return {
    authType,
    intentType,
    fileName: name,
    intentClassName: Case.pascal(name)
  }
}
