import * as fs from 'fs-extra'
import * as path from 'path'
import * as inquirer from 'inquirer'
import * as os from 'os'
import { exec } from 'child_process'

import NewCommand, {
  authTypes,
  defineLocationPath,
  askForAuthType,
  selectFolder,
  cloneRepository
} from '../../src/commands/new'
import { readFile as _readFile, ARTIFACT_FOLDER, fixturesPath } from '../helpers/utils'
import { ensureFolderExist } from '../helpers/setup'

const destination = path.join(__dirname, '..', '..', '.bearer/init')
const destinationWithViews = path.join(__dirname, '..', '..', '.bearer/initWithView')

const AUTHCONFIG = 'auth.config.json'
const PACKAGE_JSON = 'package.json'

jest.mock('inquirer')
jest.mock('child_process')

describe.each(Object.keys(authTypes))('Authentication: %s', auth => {
  const inputSet: [string, string[], string][] = [
    ['No Views', [], destination],
    ['With Views', ['--withViews'], destinationWithViews]
  ]
  inputSet.forEach(([title, args, destinationPath]) => {
    const out = path.join(destinationPath, 'new', auth)

    beforeAll(() => {
      ensureFolderExist(out)
    })

    describe(title, () => {
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

      beforeAll(() => {
        ensureFolderExist(ARTIFACT_FOLDER)
      })

      beforeEach(() => {
        resetInquirerMock()
        if (!fs.existsSync(destination)) {
          fs.mkdirSync(destination)
        }
      })

      it('asks for override and returns path if confirmed', async () => {
        expect.assertions(1)
        mockInquirer(async () => ({ override: true }))

        const location = await defineLocationPath(logger, { name, cwd: ARTIFACT_FOLDER })

        expect(location).toEqual(`${ARTIFACT_FOLDER}/${name}`)
      })

      it('prompts for foldername if not confirmed', async () => {
        const newName = 'newName'
        mockInquirer(async () => ({ override: false, response: newName }))
        expect.assertions(1)

        const location = await defineLocationPath(logger, { name, cwd: ARTIFACT_FOLDER })

        expect(location).toEqual(`${ARTIFACT_FOLDER}/${newName}`)
      })
    })
  })
})

describe('askForAuthType', () => {
  it('display all autentication types', async () => {
    const inquirerReturn = jest.fn(() => {
      return Promise.resolve({ authenticationType: 'ok' })
    })
    mockInquirer(inquirerReturn)

    await askForAuthType()

    expect(inquirerReturn.mock.calls[0]).toMatchSnapshot()
  })
})

describe('selectFolder', () => {
  describe('without any valid integration', () => {
    it('throws an no integration found error', async () => {
      await expect(selectFolder(fixturesPath('new/template-without-integration'))).rejects.toThrow(
        'No valid integration found within the cloned archive: location'
      )
    })
  })

  describe('with 1 valid integration found', () => {
    it('returns early', async () => {
      expect.assertions(1)

      const folder = await selectFolder(fixturesPath('new', 'template-with-one-integration'))

      expect(folder).toEqual({ selected: '' })
    })

    it('returns early (nested)', async () => {
      expect.assertions(1)

      const folder = await selectFolder(fixturesPath('new', 'template-with-one-integration-nested'))

      expect(folder).toEqual({ selected: 'nested-here' })
    })
  })

  describe('multiple integrations found', () => {
    it('prompt to choose one from a list', async () => {
      expect.assertions(1)
      const inquirerReturn = jest.fn(() => {
        return Promise.resolve({ selected: 'ok' })
      })
      mockInquirer(inquirerReturn)

      await selectFolder(fixturesPath('new/template-with-multiple-integrations'))

      expect(inquirerReturn.mock.calls[0]).toMatchSnapshot()
    })
  })
})

describe.only('cloneRepository', () => {
  describe('wihtout git command', () => {
    it('raise an error when git is not available', async () => {
      const logger = { debug: jest.fn(), error: jest.fn() } as any
      // @ts-ignore
      exec.mockReset()
      // @ts-ignore
      exec.mockImplementation((a, cb) => {
        return cb(new Error('ok'))
      })
      await expect(cloneRepository('ok', os.tmpdir(), logger)).rejects.toThrow(
        'git command not found in your path, please install it'
      )
    })
  })

  describe('with git installed', () => {
    beforeEach(() => {
      // @ts-ignore
      exec.mockReset()
      mockExecOnce(null, '1.0.0')
    })

    describe('no cloning errors', () => {
      it('resolves', async () => {
        expect.assertions(3)
        const logger = { debug: jest.fn(), error: jest.fn() } as any

        mockExecOnce(null, 'successfuly cloned')

        await expect(cloneRepository('ok', os.tmpdir(), logger)).resolves.toEqual(undefined)

        expect(logger.debug.mock.calls[1][0]).toContain('Running git clone ok ')
        expect(logger.error).not.toBeCalled()
      })
    })

    describe('error while cloning', () => {
      it('logs error', async () => {
        expect.assertions(1)
        const logger = { debug: jest.fn(), error: jest.fn() } as any
        mockExecOnce(new Error('failed'), '', '')

        await expect(cloneRepository('ok', os.tmpdir(), logger)).rejects.toThrow('Error while cloning the repository')
      })
    })
  })
})

function mockExecOnce(...args: any[]) {
  // @ts-ignore
  exec.mockImplementationOnce((a, cb) => {
    return cb(...args)
  })
}
function mockInquirer(response: (...args: any[]) => Promise<any>) {
  // @ts-ignore
  inquirer.prompt.mockImplementation(response)
}

function resetInquirerMock() {
  // @ts-ignore
  inquirer.prompt.mockReset()
}
