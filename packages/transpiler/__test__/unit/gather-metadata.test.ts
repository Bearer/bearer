import * as path from 'path'
import * as ts from 'typescript'

import Metadata from '../../src/metadata'
import GatherMetadata from '../../src/transformers/gather-metadata'
const fixtures = path.join(__dirname, '..', '__fixtures__')

describe('GaterMetadata transformer', () => {
  describe('a simple root component', () => {
    it('has empty input and output', () => {
      const metadata = new Metadata()
      transpileFile('root-component.tsx', metadata)
      expect(metadata).toMatchSnapshot()
    })
  })

  describe('a simple root component', () => {
    it('has input and output matching props', () => {
      const metadata = new Metadata()
      transpileFile('root-component-with-api.tsx', metadata)
      expect(metadata).toMatchSnapshot()
    })
  })
})

function transpileFile(filename: string, metadata: Metadata) {
  const files = [path.join(fixtures, 'gather-metadata', filename)]
  const compilerHost = ts.createCompilerHost({})

  const program = ts.createProgram(
    files,
    { experimentalDecorators: true, outDir: path.join(__dirname, '../../.build-transformers') },
    compilerHost
  )
  program.emit(undefined, undefined, undefined, undefined, {
    before: [GatherMetadata({ metadata })]
  })
}
