import * as Listr from 'listr'
import * as util from 'util'

const exec = util.promisify(require('child_process').exec)

export default ({
  name,
  cwd,
  command,
  env
}: {
  name: string
  cwd: string
  command: string
  env: { [key: string]: any }
}): Listr.ListrTask => ({
  title: name || `Running: ${command}`,
  task: (_ctx: any, _task: any) =>
    new Listr([
      {
        title: 'npm/yarn lookup',
        task: async (ctx: any, _task: any) =>
          exec('yarn -v')
            .then(() => {
              ctx.yarn = true
            })
            .catch(() => {
              ctx.yarn = false
            })
      },
      {
        title: 'Running command with Yarn',
        enabled: (ctx: any) => ctx.yarn === true,
        task: (_ctx: any, _task: any) => exec(`yarn ${command}`, { cwd, env })
      },
      {
        title: 'Running command with Npm',
        enabled: (ctx: any) => ctx.yarn === false,
        task: () => exec(`npm run ${command}`, { cwd, env })
      }
    ])
})
