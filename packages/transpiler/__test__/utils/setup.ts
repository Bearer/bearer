import * as fs from 'fs'
import * as globby from 'globby'
import * as path from 'path'

const buildFolder = path.join(__dirname, '../..', '.build/src')

process.env.BEARER_SCENARIO_ID = 'SPONGE_BOB'

globby.sync(path.join(buildFolder, '/**/*.tsx')).forEach(file => {
  if (fs.existsSync(file)) {
    fs.unlinkSync(file)
  }
})
