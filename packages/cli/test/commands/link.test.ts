import LinkCommand from '../../src/commands/link'
import { ensureBearerStructure } from '../helpers/setup'
import { readFile } from '../helpers/utils'

describe('Link', () => {
  let bearerPath: string

  let result: string[]

  beforeEach(() => {
    result = []
    jest.spyOn(process.stdout, 'write').mockImplementation(val => result.push(val))
    bearerPath = ensureBearerStructure()
  })

  it('does not fail :-P ', async () => {
    await LinkCommand.run(['123-integration-id', '--path', bearerPath])
    expect(result.join()).toContain('Integration successfully linked')
    expect(readFile(bearerPath, '.integrationrc')).toMatchSnapshot()
  })
})
