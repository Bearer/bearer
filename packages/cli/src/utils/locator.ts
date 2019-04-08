import * as path from 'path'

import { Config } from './setup-config'

export default class LocationProvider {
  bearerDir = ''
  integrationRoot = ''
  integrationRc!: string

  constructor(private readonly config: Config) {
    this.integrationRc = this.config.integrationConfig.config
    if (this.integrationRc) {
      this.integrationRoot = path.dirname(this.integrationRc)
      this.bearerDir = path.join(this.integrationRoot, '.bearer')
    }
  }

  integrationRootResourcePath(filename: string): string {
    return path.join(this.integrationRoot, filename)
  }

  // ~/views
  get srcViewsDir(): string {
    return path.join(this.integrationRoot, 'views')
  }

  srcViewsDirResource(name: string): string {
    return path.join(this.srcViewsDir, name)
  }

  // ~/functions
  get srcFunctionsDir(): string {
    return path.join(this.integrationRoot, 'functions')
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

  // ~/.bearer/tmp
  get buildTmpDir(): string {
    return path.join(this.bearerDir, 'tmp')
  }

  // ~/.bearer/functions
  get buildFunctionsDir(): string {
    return path.join(this.bearerDir, 'functions')
  }

  buildFunctionsResourcePath(resource: string): string {
    return path.join(this.buildFunctionsDir, resource)
  }

  get buildArtifactDir(): string {
    return path.join(this.bearerDir, 'artifacts')
  }

  buildArtifactResourcePath(resource: string): string {
    return path.join(this.buildArtifactDir, resource)
  }

  get authConfigPath(): string {
    return this.integrationRootResourcePath('auth.config.json')
  }

  get localConfigPath(): string {
    return this.integrationRootResourcePath('local.config.jsonc')
  }
}
