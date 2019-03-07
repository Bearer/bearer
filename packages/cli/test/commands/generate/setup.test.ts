import * as fs from 'fs-extra'
import * as path from 'path'

import GenerateSetup from '../../../src/commands/generate/setup'
import { ensureBearerStructure } from '../../helpers/setup'
import { readFile } from '../../helpers/utils'

const testCases: any = [
  [
    'BASIC',
    {
      authType: 'BASIC',
      setupViews: [
        { type: 'text', label: 'Username', controlName: 'username' },
        { type: 'password', label: 'Password', controlName: 'password' }
      ]
    }
  ],
  [
    'OAUTH2',
    {
      authType: 'OAUTH2',
      setupViews: [
        { type: 'text', label: 'Client ID', controlName: 'clientID' },
        { type: 'password', label: 'Client Secret', controlName: 'clientSecret' }
      ]
    }
  ],

  [
    'APIKEY',
    {
      authType: 'APIKEY',
      setupViews: [{ type: 'password', label: 'Api Key', controlName: 'apiKey' }]
    }
  ]
]

describe('generate:setup', () => {
  describe.each(testCases)('%s', (authType, auth) => {
    test('creates setup files', async () => {
      const bearerPath: string = ensureBearerStructure({
        authConfig: auth,
        folderName: authType
      })

      await GenerateSetup.run(['--force', '--path', bearerPath])

      expect(readFile(bearerPath, 'views', 'setup-action.tsx')).toMatchSnapshot()
      expect(readFile(bearerPath, 'views', 'setup-view.tsx')).toMatchSnapshot()

      expect(readFile(bearerPath, 'intents', 'saveSetup.ts')).toMatchSnapshot()
      expect(readFile(bearerPath, 'intents', 'retrieveSetup.ts')).toMatchSnapshot()

      await setTimeout(() => {}, 100)
    })
  })
})

describe('No auth', () => {
  it('create setup files ', async () => {
    const bearerPath: string = ensureBearerStructure({ authConfig: { authType: 'NONE' }, folderName: 'none' })
    await GenerateSetup.run(['--path', bearerPath])
    expect(fs.existsSync(path.join(bearerPath, 'views', 'setup-action.tsx'))).toBeFalsy()
    expect(fs.existsSync(path.join(bearerPath, 'views', 'setup-display.tsx'))).toBeFalsy()
    expect(fs.existsSync(path.join(bearerPath, 'intents', 'saveSetup.ts'))).toBeFalsy()
    expect(fs.existsSync(path.join(bearerPath, 'intents', 'retrieveSetup.ts'))).toBeFalsy()
    await setTimeout(() => {}, 100)
  })
})
