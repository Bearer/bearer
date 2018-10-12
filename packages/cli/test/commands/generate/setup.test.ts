import * as fs from 'fs-extra'
import * as path from 'path'

import GenerateSetup from '../../../src/commands/generate/setup'
import { ensureBearerStructure } from '../../helpers/setup'
import { readFile } from '../../helpers/utils'

describe('generate:setup', () => {
  describe('BASIC', () => {
    const auth = {
      authType: 'BASIC',
      setupViews: [
        { type: 'text', label: 'Username', controlName: 'username' },
        { type: 'password', label: 'Password', controlName: 'password' }
      ]
    }

    it('create setup files ', async () => {
      const result: string[] = []
      jest.spyOn(process.stdout, 'write').mockImplementation(val => {
        result.push(val)
      })

      const bearerPath: string = ensureBearerStructure({ authConfig: auth, folderName: 'BASIC' })

      await GenerateSetup.run(['--force', '--path', bearerPath])

      expect(readFile(bearerPath, 'views', 'setup-action.tsx')).toMatchSnapshot()
      expect(readFile(bearerPath, 'views', 'setup-display.tsx')).toMatchSnapshot()
      expect(result.join()).toContain('Setup components successfully generated!')
      await setTimeout(() => {}, 100)
    })
  })

  describe('OAUTH2', () => {
    const auth = {
      authType: 'OAUTH2',
      setupViews: [
        { type: 'text', label: 'Client ID', controlName: 'clientID' },
        { type: 'password', label: 'Client Secret', controlName: 'clientSecret' }
      ]
    }

    it('create setup files ', async () => {
      const result: string[] = []
      jest.spyOn(process.stdout, 'write').mockImplementation(val => {
        result.push(val)
      })
      const bearerPath: string = ensureBearerStructure({ authConfig: auth, folderName: 'OAUTH2' })

      await GenerateSetup.run(['--force', '--path', bearerPath])
      expect(readFile(bearerPath, 'views', 'setup-action.tsx')).toMatchSnapshot()
      expect(readFile(bearerPath, 'views', 'setup-display.tsx')).toMatchSnapshot()
      expect(result.join()).toContain('Setup components successfully generated!')
      await setTimeout(() => {}, 100)
    })
  })

  describe('APIKEY', () => {
    const auth = {
      authType: 'APIKEY',
      setupViews: [{ type: 'password', label: 'Api Key', controlName: 'apiKey' }]
    }

    it('create setup files ', async () => {
      const result: string[] = []
      jest.spyOn(process.stdout, 'write').mockImplementation(val => {
        result.push(val)
      })
      const bearerPath: string = ensureBearerStructure({ authConfig: auth, folderName: 'APIKEY' })

      await GenerateSetup.run(['--force', '--path', bearerPath])
      expect(readFile(bearerPath, 'views', 'setup-action.tsx')).toMatchSnapshot()
      expect(readFile(bearerPath, 'views', 'setup-display.tsx')).toMatchSnapshot()
      expect(result.join()).toContain('Setup components successfully generated!')
      await setTimeout(() => {}, 100)
    })
  })

  describe('No auth', () => {
    it('create setup files ', async () => {
      const bearerPath: string = ensureBearerStructure({ authConfig: { authType: 'NONE' }, folderName: 'none' })
      await GenerateSetup.run(['--path', bearerPath])
      expect(fs.existsSync(path.join(bearerPath, 'views', 'setup-action.tsx'))).toBeFalsy()
      expect(fs.existsSync(path.join(bearerPath, 'views', 'setup-display.tsx'))).toBeFalsy()
      await setTimeout(() => {}, 100)
    })
  })
})
