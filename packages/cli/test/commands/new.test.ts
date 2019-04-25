import * as fs from 'fs-extra'
import * as path from 'path'
import * as inquirer from 'inquirer'
import NewCommand, { authTypes, defineLocationPath } from '../../src/commands/new'
import { readFile as _readFile, ARTIFACT_FOLDER } from '../helpers/utils'

const destination = path.join(__dirname, '..', '..', '.bearer/init')
const destinationWithViews = path.join(__dirname, '..', '..', '.bearer/initWithView')

const AUTHCONFIG = 'auth.config.json'
const PACKAGE_JSON = 'package.json'

jest.mock('inquirer')
beforeAll(() => {
  if (!fs.existsSync(ARTIFACT_FOLDER)) {
    fs.mkdirSync(ARTIFACT_FOLDER)
  } else {
    fs.emptyDirSync(ARTIFACT_FOLDER)
  }
})

afterAll(() => {
  jest.unmock('inquirer')
})

describe.each(Object.keys(authTypes))('Authentication: %s', auth => {
  const inputSet: [string, string[], string][] = [
    ['No Views', [], destination],
    ['With Views', ['--withViews'], destinationWithViews]
  ]
  inputSet.forEach(([title, args, destinationPath]) => {
    const result: string[] = []
    const out = path.join(destinationPath, 'new', auth)
    describe(title, () => {
      describe(auth, () => {
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

describe('existing path', () => {
  it('asks for confirmation', () => {})

  describe('when OK', () => {
    it('generates integration', () => {})
  })

  describe('when canceled', () => {
    it('stops and return', () => {})
  })
})

describe('defineLocationPath', () => {
  const logger = { warn: jest.fn() }

  describe('with force', () => {
    it('returns the path + the name', async () => {
      expect.assertions(1)

      const location = await defineLocationPath(logger, { name: 'whatever', cwd: '/tmp/ok', force: true })

      expect(location).toEqual('/tmp/ok/whatever')
    })
  })

  describe('without force', () => {
    describe('when folder does not exist', () => {
      it('returns path', async () => {
        const name = Date.now().toString()
        expect.assertions(1)

        const location = await defineLocationPath(logger, { name, cwd: '/tmp/ok' })

        expect(location).toEqual(`/tmp/ok/${name}`)
      })
    })

    describe('when folder exists', () => {
      const name = 'existing'
      const destination = path.join(ARTIFACT_FOLDER, name)

      beforeEach(() => {
        resetInquirerMock()
        if (!fs.existsSync(destination)) {
          fs.mkdirSync(destination)
        }
      })

      it('asks for override and returns path if confirmed', async () => {
        expect.assertions(1)
        mockInquirer({ override: true })

        const location = await defineLocationPath(logger, { name, cwd: ARTIFACT_FOLDER })

        expect(location).toEqual(`${ARTIFACT_FOLDER}/${name}`)
      })

      it('prompts for foldername if not confirmed', async () => {
        const newName = 'newName'
        mockInquirer({ override: false, response: newName })
        expect.assertions(1)

        const location = await defineLocationPath(logger, { name, cwd: ARTIFACT_FOLDER })

        expect(location).toEqual(`${ARTIFACT_FOLDER}/${newName}`)
      })
    })
  })
})

function mockInquirer(response: any) {
  // @ts-ignore
  inquirer.prompt.mockImplementation(() => {
    return Promise.resolve(response)
  })
}

function resetInquirerMock() {
  // @ts-ignore
  inquirer.prompt.mockReset()
}
