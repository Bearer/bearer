import * as Listr from 'listr'
import * as util from 'util'

const exec = util.promisify(require('child_process').exec)

export default ({ cwd }: { cwd: string }): Listr.ListrTask => ({
  title: 'Install dependencies',
  task: async (_ctx: any, _task: any) =>
    new Listr([
      {
        title: 'npm/yarn lookup',
        task: async (ctx: any, task: any) =>
          exec('yarn -v')
            .then(() => {
              ctx.yarn = true
            })
            .catch(() => {
              ctx.yarn = false
              task.skip('Yarn not available, install it via `npm install -g yarn`')
            })
      },
      {
        title: 'Installing scenario dependencies with yarn',
        enabled: (ctx: any) => ctx.yarn === true,
        task: async (_ctx: any, _task: any) => exec('yarn install', { cwd })
      },
      {
        title: 'Installing scenario dependencies with npm',
        enabled: (ctx: any) => ctx.yarn === false,
        task: async () => exec('yarn install', { cwd })
      }
    ])
})
