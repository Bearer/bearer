import * as fs from 'fs-extra'
import * as path from 'path'

import GenerateSetup from '../../../src/commands/generate/setup'
import { ensureBearerStructure } from '../../helpers/setup'
import { readFile } from '../../helpers/utils'

const testCases: any = [
  ['BASIC', { authType: 'BASIC' }],
  ['OAUTH1', { authType: 'OAUTH1' }],
  ['OAUTH2', { authType: 'OAUTH2' }],
  ['APIKEY', { authType: 'APIKEY' }]
]

describe('generate:setup', () => {
  describe.each(testCases)('%s', (authType, auth) => {
    test('creates setup files', async () => {
      const bearerPath: string = ensureBearerStructure({
        clean: true,
        authConfig: auth,
        folderName: authType,
        withViews: true
      })
      await GenerateSetup.run(['--force', '--path', bearerPath])

      expect(readFile(bearerPath, 'views', 'setup-action.tsx')).toMatchSnapshot()
      expect(readFile(bearerPath, 'views', 'setup-view.tsx')).toMatchSnapshot()

      expect(readFile(bearerPath, 'functions', 'saveSetup.ts')).toMatchSnapshot()
      expect(readFile(bearerPath, 'functions', 'retrieveSetup.ts')).toMatchSnapshot()

      await setTimeout(() => {}, 100)
    })
  })
})

describe('No Auth', () => {
  it('create setup files ', async () => {
    const bearerPath: string = ensureBearerStructure({ authConfig: { authType: 'NONE' }, folderName: 'none' })
    await GenerateSetup.run(['--path', bearerPath])
    expect(fs.existsSync(path.join(bearerPath, 'views', 'setup-action.tsx'))).toBeFalsy()
    expect(fs.existsSync(path.join(bearerPath, 'views', 'setup-display.tsx'))).toBeFalsy()
    expect(fs.existsSync(path.join(bearerPath, 'functions', 'saveSetup.ts'))).toBeFalsy()
    expect(fs.existsSync(path.join(bearerPath, 'functions', 'retrieveSetup.ts'))).toBeFalsy()
    await setTimeout(() => {}, 100)
  })
})

describe('Custom Auth', () => {
  it('create setup files ', async () => {
    const bearerPath: string = ensureBearerStructure({ authConfig: { authType: 'CUSTOM' }, folderName: 'none' })
    await GenerateSetup.run(['--path', bearerPath])
    expect(fs.existsSync(path.join(bearerPath, 'views', 'setup-action.tsx'))).toBeFalsy()
    expect(fs.existsSync(path.join(bearerPath, 'views', 'setup-display.tsx'))).toBeFalsy()
    expect(fs.existsSync(path.join(bearerPath, 'functions', 'saveSetup.ts'))).toBeFalsy()
    expect(fs.existsSync(path.join(bearerPath, 'functions', 'retrieveSetup.ts'))).toBeFalsy()
    await setTimeout(() => {}, 100)
  })
})
