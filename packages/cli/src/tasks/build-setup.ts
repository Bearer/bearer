import * as Listr from 'listr'
import GenerateSetup from '../commands/generate/setup'
import { copyFiles } from '../utils/helpers'

export default ({ cmd, vars }: { cmd: GenerateSetup; vars: any }): Listr.ListrTask[] => [
  {
    title: 'Generating setup components',
    task: async () => {
      await copyFiles(cmd, 'generate/setup/components', cmd.locator.srcViewsDir, vars)
      return true
    }
  },
  {
    title: 'Generating setup intents',
    task: async () => {
      await copyFiles(cmd, `generate/setup/intents/${cmd.scenarioAuthConfig.authType}`, cmd.locator.srcIntentsDir, vars)
      return true
    }
  }
]
