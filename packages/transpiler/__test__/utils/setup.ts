import * as fs from 'fs'
import * as globby from 'globby'
import * as path from 'path'

const buildFolder = path.join(__dirname, '..', '.build/src')

process.env.BEARER_SCENARIO_ID = 'SPONGE_BOB'

globby.sync(['**/*.tsx', '**/*.ts', '**/*.json'], { cwd: buildFolder }).forEach(file => {
  const filePath = path.join(buildFolder, file)
  if (fs.existsSync(filePath)) {
    fs.unlinkSync(filePath)
  }
})
