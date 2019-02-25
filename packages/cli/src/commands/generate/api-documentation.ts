import { flags } from '@oclif/command'
import * as fs from 'fs-extra'
import * as path from 'path'

import BaseCommand from '../../base-command'
import { RequireScenarioFolder } from '../../utils/decorators'
import { OpenApiSpecGenerator } from '../../utils/generators'

const OPEN_API_SPEC = 'openapi.json'

export default class GenerateApiDocumentation extends BaseCommand {
  static description = 'Generate an openapi REST documentation'
  static hidden = true
  static flags = {
    ...BaseCommand.flags,
    soft: flags.boolean({})
  }

  static args = []

  @RequireScenarioFolder()
  async run() {
    try {
      const { flags } = this.parse(GenerateApiDocumentation)
      const { srcIntentsDir, buildViewsComponentsDir } = this.locator

      const { scenarioTitle, scenarioUuid } = this.bearerConfig
      const spec = flags.soft
        ? {}
        : await new OpenApiSpecGenerator(srcIntentsDir, { scenarioTitle, scenarioUuid }).build()

      fs.ensureDirSync(buildViewsComponentsDir)

      const openApiSpecPath = path.join(buildViewsComponentsDir, OPEN_API_SPEC)
      const fileExists = fs.existsSync(openApiSpecPath)
      await fs.writeFile(openApiSpecPath, spec)
      const action = fileExists ? 'updated' : 'generated'
      this.success(`File ${OPEN_API_SPEC} ${action}!`)
    } catch (e) {
      this.error(e)
    }
  }
}
