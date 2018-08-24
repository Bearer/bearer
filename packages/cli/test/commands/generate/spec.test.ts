import GenerateSpec from '../../../src/commands/generate/spec'
import { ensureBearerStructure } from '../../helpers/setup'
import { readFile } from '../../helpers/utils'

describe('Generate', () => {
  let bearerPath: string
  let result: Array<string>

  beforeEach(() => {
    result = []
    jest.spyOn(process.stdout, 'write').mockImplementation(val => result.push(val))
    bearerPath = ensureBearerStructure()
  })

  describe('generate:spec', () => {
    it('creates a spce file', async () => {
      await GenerateSpec.run(['--force', '--path', bearerPath])
      expect(result.join()).toContain('Spec file successfully generated!')
      expect(readFile(bearerPath, 'spec.ts')).toMatchSnapshot()
    })
  })
})
