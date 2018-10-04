import * as path from 'path'
import * as ts from 'typescript'

import GatherMetadata from '../../src/transformers/gather-metadata'
import { Metadata } from '../../src/types'
const fixtures = path.join(__dirname, '..', '__fixtures__')

describe('GaterMetadata transformer', () => {
  describe('a simple root component', () => {
    it('has empty input and output', () => {
      const metadata = { components: [] }
      transpileFile('root-component.tsx', metadata)
      expect(metadata).toMatchSnapshot()
    })
  })

  describe('a simple root component', () => {
    it('has input and output matching props', () => {
      const metadata = { components: [] }
      transpileFile('root-component-with-api.tsx', metadata)
      expect(metadata).toMatchSnapshot()
    })
  })
})

function transpileFile(filename: string, metadata: Metadata) {
  const compilerHost = ts.createCompilerHost({})
  const program = ts.createProgram(
    [path.join(fixtures, 'gather-metadata', filename)],
    { experimentalDecorators: true, outDir: path.join(__dirname, '../../.build-transformers') },
    compilerHost
  )
  program.emit(undefined, undefined, undefined, undefined, {
    before: [GatherMetadata({ metadata })]
  })
}
