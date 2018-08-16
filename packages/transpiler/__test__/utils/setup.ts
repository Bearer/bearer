import fs from 'fs'
import globby from 'globby'
import path from 'path'

const buildFolder = path.join(__dirname, '../..', '.build/src')

process.env.BEARER_SCENARIO_ID = 'SPONGE_BOB'

globby.sync(path.join(buildFolder, '/**/*.tsx')).forEach(file => {
  fs.unlinkSync(file)
})
