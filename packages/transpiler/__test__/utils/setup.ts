import fs from 'fs'
import path from 'path'
import { TranspilerFactory } from './transpiler'

const fixtures = path.join(__dirname, '__fixtures__')
const buildFolder = path.join(__dirname, '../..', '.build/src')

console.log('[BEARER]', 'cleaning build folder', buildFolder)

process.env.BEARER_SCENARIO_ID = 'SPONGE_BOB'

fs.readdirSync(buildFolder).forEach(file => {
  if (file.match(/\.tsx?$/)) {
    console.log('removing', file)
    fs.unlinkSync(path.join(buildFolder, file))
  }
})

const transpiler = TranspilerFactory({
  ROOT_DIRECTORY: fixtures,
  srcFolder: '../../__fixtures__'
})
transpiler.run()
