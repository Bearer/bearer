import * as fs from 'fs-extra'
import * as path from 'path'

type TSetupConfig = {
  clean?: boolean
  authConfig?: any
  folderName?: string
  withViews?: boolean
}

export function ensureBearerStructure({
  clean = true,
  authConfig,
  folderName = 'fakeintegration',
  withViews = false
}: TSetupConfig = {}): string {
  const bearerFolder = path.join(__dirname, '..', '..', '.bearer', folderName)
  if (!fs.existsSync(bearerFolder)) {
    fs.mkdirpSync(bearerFolder)
  }

  if (!fs.existsSync(bearerFolder)) {
    fs.mkdirpSync(bearerFolder)
  }
  if (clean) {
    fs.emptyDirSync(bearerFolder)
  }
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
        authType: 'APIKEY',
        setupViews: [
          {
            type: 'password',
            label: 'Api Key',
            controlName: 'apiKey'
          }
        ]
      }
    )
  )
  return bearerFolder
}
