import * as path from 'path'
import * as ts from 'typescript'
import * as TJS from 'typescript-json-schema'

import Metadata from '../../src/metadata'
import GatherIO from '../../src/transformers/gather-input-output'
const fixtures = path.join(__dirname, '..', '__fixtures__')

describe('GatherIO transformer', () => {
  describe('a simple root component', () => {
    it('has inputs and outputs added', () => {
      const metadata = new Metadata()
      metadata.registerComponent({
        classname: 'FeatDisplayComponent',
        fileName: '/__test__/__fixtures__/gather-input-output/root-component-transformed.tsx',
        finalTagName: 'complex-feature-display',
        name: 'complex',
        initialTagName: 'feat-display',
        isRoot: true
      })

      transpileFile('root-component-transformed.tsx', metadata)
      expect(metadata).toMatchSnapshot()
    })
  })
})

function transpileFile(filename: string, metadata: Metadata) {
  const files = [path.join(fixtures, 'gather-input-output', filename)]
  const config = ts.readConfigFile(
    path.join(__dirname, '../../../cli/templates/start', 'tsconfig.json'),
    ts.sys.readFile
  )
  const compilerHost = ts.createCompilerHost({})
  const programGenerator = TJS.getProgramFromFiles(files, config.config.compilerOptions, './ok')

  const generator = TJS.buildGenerator(programGenerator, {
    required: true
  })
  const program = ts.createProgram(
    files,
    { experimentalDecorators: true, outDir: path.join(__dirname, '../../.build-transformers') },
    compilerHost
  )
  program.emit(undefined, undefined, undefined, undefined, {
    before: [GatherIO({ metadata, generator })]
  })
}
