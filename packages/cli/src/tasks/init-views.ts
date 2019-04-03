import * as Listr from 'listr'
import * as fs from 'fs'
import * as path from 'path'
// @ts-ignore
import * as set from 'lodash.set'

import Command from '../commands/new'
import { copyFiles } from '../utils/helpers'

export default ({
  cmd,
  vars
}: {
  cmd: Command
  vars: {
    bearerTagVersion: string
    componentName: string
    setup: string
  }
}): Listr.ListrTask => ({
  title: 'Generating views directory',
  task: async (ctx: any) => {
    const files = await copyFiles(cmd, 'init/views', cmd.copyDestFolder, vars, true)
    const packageJson = path.join(cmd.copyDestFolder, 'package.json')
    const projectPackage: any = JSON.parse(fs.readFileSync(packageJson, { encoding: 'utf8' }))

    // inject dependencies
    set(projectPackage, 'dependencies.@bearer/core', vars.bearerTagVersion)
    set(projectPackage, 'dependencies.@bearer/ui', vars.bearerTagVersion)
    set(projectPackage, 'devDependencies.@bearer/js', vars.bearerTagVersion)
    set(projectPackage, 'peerDependencies.@bearer/js', vars.bearerTagVersion)

    fs.writeFileSync(packageJson, JSON.stringify(projectPackage, null, 2))

    ctx.files = [...ctx.files, ...files]
    return true
  }
})
