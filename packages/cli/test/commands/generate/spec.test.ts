import GenerateSpec from '../../../src/commands/generate/spec'
import { ensureBearerStructure } from '../../helpers/setup'
import { readFile } from '../../helpers/utils'

describe('Generate', () => {
  let bearerPath: string
  let result: string[]

  beforeEach(() => {
    result = []
    jest.spyOn(process.stdout, 'write').mockImplementation(val => {
      result.push(val)
      return true
    })
    bearerPath = ensureBearerStructure()
  })

  describe('generate:spec', () => {
    xit('creates a spec file', async () => {
      await GenerateSpec.run(['--force', '--path', bearerPath])
      expect(result.join()).toContain('Spec file successfully generated!')
      expect(readFile(bearerPath, 'spec.ts')).toMatchSnapshot()
    })
  })
})
