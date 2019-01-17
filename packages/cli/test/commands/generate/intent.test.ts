import GenerateIntent from '../../../src/commands/generate/intent'
import { ensureBearerStructure } from '../../helpers/setup'
import { readFile } from '../../helpers/utils'
import { Authentications } from '@bearer/types/lib/authentications'

describe('Generate', () => {
  let bearerPath: string
  let result: string[]

  describe.each(Object.values(Authentications))(`%s - generate:intent`, authType => {
    beforeAll(() => {
      result = []
      jest.spyOn(process.stdout, 'write').mockImplementation(val => result.push(val))
      bearerPath = ensureBearerStructure({ clean: true, authConfig: { authType }, folderName: authType })
    })

    it('Fetch intent', async () => {
      await GenerateIntent.run(['FetchDataIntent', '-t', 'fetch', '--path', bearerPath])
      expect(result.join()).toContain('Intent generated')
      expect(readFile(bearerPath, 'intents', 'FetchDataIntent.ts')).toMatchSnapshot()
    })

    it('Save intent', async () => {
      await GenerateIntent.run(['SaveIntent', '-t', 'save', '--path', bearerPath])
      expect(result.join()).toContain('Intent generated')
      expect(readFile(bearerPath, 'intents', 'SaveIntent.ts')).toMatchSnapshot()
    })
  })
})
