import { ScenarioConfig } from './types'
import * as path from 'path'

export default class LocationProvider {
  scenarioRoot: string
  constructor(private readonly config: ScenarioConfig) {
    this.scenarioRoot = path.dirname(this.config.scenarioConfig.config)
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
