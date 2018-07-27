import { ScenarioConfig } from './types'
import * as path from 'path'

export default class LocationProvider {
  bearerDir: string
  scenarioRoot: string
  scenarioRc: string

  constructor(private readonly config: ScenarioConfig) {
    this.scenarioRc = this.config.scenarioConfig.config
    if (this.scenarioRc) {
      this.scenarioRoot = path.dirname(this.scenarioRc)
      this.bearerDir = path.join(this.scenarioRoot, '.bearer')
    }
  }

  scenarioRootResourcePath(filename: string): string {
    return path.join(this.scenarioRoot, filename)
  }

  // ~/views
  get srcViewsDir(): string {
    return path.join(this.scenarioRoot, 'views')
  }
  // ~/intents
  get srcIntentsDir(): string {
    return path.join(this.scenarioRoot, 'intents')
  }

  buildViewsResourcePath(resource: string): string {
    return path.join(this.buildViewsDir, resource)
  }

  // ~/.bearer/views
  get buildViewsDir(): string {
    return path.join(this.bearerDir, 'views')
  }

  // ~/.bearer/views/src
  get buildViewsComponentsDir(): string {
    return path.join(this.buildViewsDir, 'src')
  }

  // ~/.bearer/intents
  get buildIntentsDir(): string {
    return path.join(this.bearerDir, 'intents')
  }

  buildIntentsResourcePath(resource: string): string {
    return path.join(this.buildIntentsDir, resource)
  }

  get intentsArtifactDir(): string {
    return path.join(this.bearerDir, 'artifacts')
  }

  intentsArtifactResourcePath(resource: string): string {
    return path.join(this.intentsArtifactDir, resource)
  }

  get authConfigPath(): string {
    return this.scenarioRootResourcePath('auth.config.json')
  }
}
