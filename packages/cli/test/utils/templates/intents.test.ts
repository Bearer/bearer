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

  describe.each(Object.values(Authentications))('%s', (auth: Authentications) => {
    describe.each(Object.values(IntentType))('%s', (intentType: IntentType) => {
      it('generate valide TS file', async () => {
        const command = { silent: true, locator: { srcIntentsDir: destination } }
        const files = await generateIntent(command as any, auth, intentType, `${auth}-${intentType}-Intent`)
        const options = ts.convertCompilerOptionsFromJson(compilerOptions, 'ok')
        const program = ts.createProgram(files, { ...options.options, noEmit: true });
        const emitResult = program.emit()
        const allDiagnostics = ts.getPreEmitDiagnostics(program).concat(emitResult.diagnostics)
        expect(allDiagnostics).toHaveLength(0)
      })
    })
  })
})