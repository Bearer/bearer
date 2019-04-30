import * as ts from 'typescript'
import * as fs from 'fs-extra'

import { Authentications } from '@bearer/types/lib/authentications'
import FunctionType from '@bearer/types/lib/function-types'

import generateFunction from '../../../src/utils/templates/functions'
import compilerOptions from '../../../src/utils/function-ts-compiler-options'
import { artifactPath } from '../../helpers/utils'
const destination = artifactPath('generated-functions')

describe('functions generator', () => {
  beforeAll(() => {
    if (fs.existsSync(destination)) {
      fs.emptyDirSync(destination)
    }
  })

  describe.each(Object.values(Authentications))('When %s', (auth: Authentications) => {
    describe.each(Object.values(FunctionType))('function type: %s', (functionType: FunctionType) => {
      let files: string[] = []
      let diagnostics: ts.Diagnostic[] = []

      beforeAll(async () => {
        const command = { silent: true, locator: { srcFunctionsDir: destination } }
        files = await generateFunction(command as any, auth, functionType, `${auth}-${functionType}-Function`)
        const options = ts.convertCompilerOptionsFromJson(compilerOptions, 'ok')
        const program = ts.createProgram(files, { ...options.options, noEmit: true })
        const emitResult = program.emit()
        diagnostics = ts.getPreEmitDiagnostics(program).concat(emitResult.diagnostics)
      })

      it('generates valide TS file', async () => {
        expect(diagnostics).toHaveLength(0)
      })

      it('matches snapshot', () => {
        expect(fs.readFileSync(files[0], { encoding: 'utf8' })).toMatchSnapshot()
      })
    })
  })
})
