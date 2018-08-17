import * as Case from 'case'
import * as copy from 'copy-template-dir'
import * as del from 'del'
import * as detect from 'detect-file'
import * as path from 'path'
import * as rc from 'rc'

import Locator from '../../locationProvider'

export function generateSetup({
  emitter,
  locator,
  deleteSetup
}: {
  emitter: any
  locator: Locator
  deleteSetup?: boolean
}) {
  return new Promise(async (resolve, reject) => {
    try {
      const authConfig = require(locator.authConfigPath)
      const scenarioConfig = rc('scenario')
      const { scenarioTitle } = scenarioConfig
      const configKey = 'setupViews'
      const inDir = path.join(__dirname, '../templates/generate/setup')
      const outDir = locator.srcViewsDir

      if (deleteSetup) {
        await del(`${outDir}/setup*.tsx`).then(paths => {
          paths.forEach(path => emitter.emit('generateTemplate:deleteFiles', path))
        })
      }

      if (isSetupExist(locator)) {
        resolve()
      } else if (authConfig[configKey] && authConfig[configKey].length) {
        const vars = {
          componentName: Case.pascal(scenarioTitle),
          componentTagName: Case.kebab(scenarioTitle),
          fields: JSON.stringify(authConfig[configKey])
        }
        copy(inDir, outDir, vars, (err, createdFiles) => {
          if (err) throw err
          createdFiles.forEach(filePath => emitter.emit('generateTemplate:fileGenerated', filePath))
          resolve()
        })
      } else {
        emitter.emit('generateTemplate:skipped', configKey)
        resolve()
      }
    } catch (error) {
      emitter.emit('generateTemplate:error', error.toString())
      reject(error)
    }
  })
}

function isSetupExist(locator): boolean {
  return (
    detect(path.join(locator.srcViewsDir, 'setup-action.tsx')) &&
    detect(path.join(locator.srcViewsDir, 'setup-display.tsx'))
  )
}
