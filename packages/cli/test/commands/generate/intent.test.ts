import GenerateIntent from '../../../src/commands/generate/intent'
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

  describe('generate:intent', () => {
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

    it('Retrieve intent', async () => {
      await GenerateIntent.run(['RetrieveIntent', '-t', 'retrieve', '--path', bearerPath])
      expect(result.join()).toContain('Intent generated')
      expect(readFile(bearerPath, 'intents', 'RetrieveIntent.ts')).toMatchSnapshot()
    })
  })
})
