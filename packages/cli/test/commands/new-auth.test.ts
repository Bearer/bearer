import * as path from 'path'

import NewCommand, { authTypes } from '../../src/commands/new'
import { readFile as _readFile, artifactPath } from '../helpers/utils'
import { ensureFolderExist } from '../helpers/setup'

const destination = artifactPath('init')
const destinationWithViews = artifactPath('initWithViews')

const AUTHCONFIG = 'auth.config.json'
const PACKAGE_JSON = 'package.json'

jest.mock('inquirer')
jest.mock('child_process')

afterAll(() => {
  jest.unmock('inquirer')
  jest.unmock('child_process')
})

describe.each(Object.keys(authTypes))('Authentication: %s', auth => {
  const inputSet: [string, string[], string][] = [
    ['No Views', [], destination],
    ['With Views', ['--withViews'], destinationWithViews]
  ]
  inputSet.forEach(([title, args, destinationPath]) => {
    describe(title, () => {
      const out = path.join(destinationPath, 'new', auth)

      beforeAll(async () => {
        await ensureFolderExist(out)
      })
      describe(auth, () => {
        const result: string[] = []
        beforeAll(async () => {
          jest.spyOn(process.stdout, 'write').mockImplementation(val => {
            result.push(val)
            return true
          })
          await NewCommand.run([`${auth}Integration`, '-a', auth, '--skipInstall', '--path', out, ...args])
        })

        it('produces an output', () => {
          const outPut = result.sort().join()
          expect(outPut).toMatchSnapshot()
        })

        it('creates an auth.config.json file', () => {
          expect(_readFile(out, `${auth}Integration`, AUTHCONFIG)).toMatchSnapshot()
        })

        it('creates a package.json file', () => {
          expect(_readFile(out, `${auth}Integration`, PACKAGE_JSON)).toMatchSnapshot()
        })
      })
    })
  })
})
