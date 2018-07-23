import { ScenarioConfig } from './types'
import * as path from 'path'

export default class LocationProvider {
  scenarioRoot: string
  scenarioRc: string

  constructor(private readonly config: ScenarioConfig) {
    this.scenarioRc = this.config.scenarioConfig.config
    if (this.scenarioRc) {
      this.scenarioRoot = path.dirname(this.scenarioRc)
    }
  }

  scenarioRootFile(filename: string): string {
    return path.join(this.scenarioRoot, filename)
  }
  // ~/screens
  get srcScreenDir(): string {
    return path.join(this.scenarioRoot, 'screens')
  }

  // ~/.build/
  get buildDir(): string {
    return path.join(this.scenarioRoot, '.build')
  }

  // ~/.build/src
  get buildScreenDir(): string {
    return path.join(this.buildDir, 'src')
  }
}
