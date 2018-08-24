import * as fs from 'fs-extra'
import * as path from 'path'

type TSetupConfig = {
  clean?: boolean
  authConfig?: any
  folderName?: string
}

export function ensureBearerStructure({
  clean = true,
  authConfig,
  folderName = 'fakescenario'
}: TSetupConfig = {}): string {
  const bearerFolder = path.join(__dirname, '..', '.bearer', folderName)
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
