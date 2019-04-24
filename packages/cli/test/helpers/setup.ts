import * as fs from 'fs-extra'
import * as path from 'path'

type TSetupConfig = {
  clean?: boolean
  authConfig?: any
  folderName?: string
  withViews?: boolean
}

export const ARTIFACT_ROOT = path.join(__dirname, '..', '..', '.bearer')

export function cleanArtifactFolder(name: string) {
  const bearerFolder = path.join(ARTIFACT_ROOT, name)
  if (!fs.existsSync(bearerFolder)) {
    fs.mkdirpSync(bearerFolder)
  }

  fs.emptyDirSync(bearerFolder)

  return bearerFolder
}

export function ensureBearerStructure({
  authConfig,
  folderName = 'fakeintegration',
  withViews = false
}: TSetupConfig = {}): string {
  const bearerFolder = cleanArtifactFolder(folderName)

  if (withViews) {
    const views = path.join(bearerFolder, 'views')
    if (!fs.existsSync(views)) {
      fs.mkdirpSync(views)
    }
  }
  fs.writeFileSync(path.join(bearerFolder, '.integrationrc'), 'integrationTile=test')
  fs.writeFileSync(
    path.join(bearerFolder, 'auth.config.json'),
    JSON.stringify(
      authConfig || {
        authType: 'APIKEY'
      }
    )
  )
  return bearerFolder
}
