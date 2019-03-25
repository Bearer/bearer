import * as Listr from 'listr'
import Command from '../base-command'
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
    const files = await copyFiles(cmd, 'init/views', cmd.locator.integrationRoot, vars)
    // TODO: fix returned path
    ctx.files = [...ctx.files, ...files]
    return true
  }
})
