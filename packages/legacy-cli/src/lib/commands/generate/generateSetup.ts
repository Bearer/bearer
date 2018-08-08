import * as rc from 'rc'
import * as path from 'path'
import * as del from 'del'
import * as Case from 'case'
import * as copy from 'copy-template-dir'
import * as fs from 'fs'
import Locator from '../../locationProvider'

export async function generateSetup({
  emitter,
  locator,
  deleteSetup
}: {
  emitter: any
  locator: Locator
  deleteSetup?: boolean
}) {
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

    if (authConfig[configKey] && authConfig[configKey].length) {
      const vars = {
        componentName: Case.pascal(scenarioTitle),
        componentTagName: Case.kebab(scenarioTitle),
        fields: JSON.stringify(authConfig[configKey])
      }
      copy(inDir, outDir, vars, (err, createdFiles) => {
        if (err) throw err
        createdFiles.forEach(filePath => emitter.emit('generateTemplate:fileGenerated', filePath))
      })
    } else {
      throw new Error('Configuration file is incorrect or missing')
    }
  } catch (error) {
    emitter.emit('generateTemplate:error', error.toString())
  }
}
