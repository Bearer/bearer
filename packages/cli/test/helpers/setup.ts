import * as fs from 'fs-extra'
import * as path from 'path'

type TSetupConfig = {
  clean?: boolean
  authConfig?: any
}

export function ensureBearerStructure({ clean = true, authConfig }: TSetupConfig = {}): string {
  const bearerFolder = path.join(__dirname, '..', '.bearer')
  if (!fs.existsSync(bearerFolder)) {
    fs.mkdirpSync(bearerFolder)
  }
  if (clean) {
    fs.emptyDirSync(bearerFolder)
  }
  fs.writeFileSync(path.join(bearerFolder, '.scenariorc'), 'scenarioTile=test')
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
