import * as ts from 'typescript'
import * as fs from 'fs-extra'

import { Authentications } from '@bearer/types/lib/authentications';
import IntentType from '@bearer/types/lib/intent-types'

import * as path from 'path'

import generateIntent from '../../../src/utils/templates/intents'
import compilerOptions from '../../../src/utils/intent-ts-compiler-options'
const destination = path.join(__dirname, '../../../.bearer/generated-intents')

describe('intents generator', () => {
  beforeAll(() => {
    if(fs.existsSync(destination)){
      fs.emptyDirSync(destination)
    }
  })

  describe.each(Object.values(Authentications))('When %s', (auth: Authentications) => {
    describe.each(Object.values(IntentType))('intent type: %s', (intentType: IntentType) => {
      let files: string[] = []
      let diagnostics: ts.Diagnostic[] = []

      beforeAll(async () => {
        const command = { silent: true, locator: { srcIntentsDir: destination } }
        files = await generateIntent(command as any, auth, intentType, `${auth}-${intentType}-Intent`)
        const options = ts.convertCompilerOptionsFromJson(compilerOptions, 'ok')
        const program = ts.createProgram(files, { ...options.options, noEmit: true });
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