import * as fs from 'fs-extra'
import * as path from 'path'

import GenerateSetup from '../../../src/commands/generate/setup'
import { ensureBearerStructure } from '../../helpers/setup'
import { readFile } from '../../helpers/utils'

describe('Generate', () => {
  let bearerPath: string
  let result: Array<string>

  const sampleConfig = [
    {
      authType: 'APIKEY',
      setupViews: [
        {
          type: 'password',
          label: 'Api Key',
          controlName: 'apiKey'
        }
      ]
    },
    {
      authType: 'OAUTH2',
      setupViews: [
        {
          type: 'text',
          label: 'Client ID',
          controlName: 'clientID'
        },
        {
          type: 'password',
          label: 'Client Secret',
          controlName: 'clientSecret'
        }
      ]
    },
    {
      authType: 'BASIC',
      setupViews: [
        {
          type: 'text',
          label: 'Username',
          controlName: 'username'
        },
        {
          type: 'password',
          label: 'Password',
          controlName: 'password'
        }
      ]
    }
  ]

  describe('generate:setup', () => {
    beforeEach(() => {
      jest.spyOn(process.stdout, 'write').mockImplementation(val => {
        result.push(val)
      })
    })
    sampleConfig.map(auth => {
      describe(auth.authType, () => {
        beforeEach(() => {
          result = []
          bearerPath = ensureBearerStructure({ authConfig: auth, folderName: auth.authType })
        })
        it('create setup files ', async () => {
          await GenerateSetup.run(['--force', '--path', bearerPath])
          expect(readFile(bearerPath, 'views', 'setup-action.tsx')).toMatchSnapshot()
          expect(readFile(bearerPath, 'views', 'setup-display.tsx')).toMatchSnapshot()
          expect(result.join()).toContain('Setup components successfully generated!')
        })
      })
    })

    describe('No auth', () => {
      it('create setup files ', async () => {
        bearerPath = ensureBearerStructure({ authConfig: { authType: 'NONE' }, folderName: 'none' })
        await GenerateSetup.run(['--path', bearerPath])
        expect(fs.existsSync(path.join(bearerPath, 'views', 'setup-action.tsx'))).toBeFalsy()
        expect(fs.existsSync(path.join(bearerPath, 'views', 'setup-display.tsx'))).toBeFalsy()
      })
    })
  })
})
