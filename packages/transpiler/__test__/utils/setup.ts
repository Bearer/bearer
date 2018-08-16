import fs from 'fs'
import path from 'path'

const buildFolder = path.join(__dirname, '../..', '.build/src')

console.log('[BEARER]', 'cleaning build folder', buildFolder)

process.env.BEARER_SCENARIO_ID = 'SPONGE_BOB'

fs.readdirSync(buildFolder).forEach(file => {
  if (file.match(/\.tsx?$/)) {
    console.log('removing', file)
    fs.unlinkSync(path.join(buildFolder, file))
  }
})
