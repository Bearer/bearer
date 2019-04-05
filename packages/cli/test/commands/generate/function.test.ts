import GenerateFunction from '../../../src/commands/generate/function'
import { ensureBearerStructure } from '../../helpers/setup'
import { readFile } from '../../helpers/utils'
import { Authentications } from '@bearer/types/lib/authentications'

describe('Generate', () => {
  let bearerPath: string
  let result: string[]

  describe.each(Object.values(Authentications))(`%s - generate:function`, authType => {
    beforeAll(() => {
      result = []
      jest.spyOn(process.stdout, 'write').mockImplementation(val => {
        result.push(val)
        return true
      })
      bearerPath = ensureBearerStructure({ authConfig: { authType }, folderName: authType })
    })

    it('Fetch function', async () => {
      await GenerateFunction.run(['FetchDataFunction', '--path', bearerPath])
      expect(result.join()).toContain('Function generated')
      expect(readFile(bearerPath, 'functions', 'FetchDataFunction.ts')).toMatchSnapshot()
    })
  })
})
