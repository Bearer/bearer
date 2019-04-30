import * as fs from 'fs-extra'
import * as path from 'path'
import * as inquirer from 'inquirer'
import * as os from 'os'
import { exec } from 'child_process'

import { defineLocationPath, askForAuthType, selectFolder, cloneRepository } from '../../src/commands/new'
import { readFile as _readFile, ARTIFACT_FOLDER, fixturesPath } from '../helpers/utils'
import { ensureFolderExist } from '../helpers/setup'

jest.mock('inquirer')
jest.mock('child_process')

afterAll(() => {
  jest.unmock('inquirer')
  jest.unmock('child_process')
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

describe.skip('defineLocationPath', () => {
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

describe.skip('askForAuthType', () => {
  it('display all autentication types', async () => {
    const inquirerReturn = jest.fn(() => {
      return Promise.resolve({ authenticationType: 'ok' })
    })
    mockInquirer(inquirerReturn)

    await askForAuthType()

    expect(inquirerReturn.mock.calls[0]).toMatchSnapshot()
  })
})

describe.skip('selectFolder', () => {
  describe('without any valid integration', () => {
    it('throws an no integration found error', async () => {
      await expect(selectFolder(fixturesPath('new/template-without-integration'), {})).rejects.toThrow(
        'No valid integrations found'
      )
    })
  })

  describe('with 1 valid integration found', () => {
    it('returns early', async () => {
      expect.assertions(1)

      const folder = await selectFolder(fixturesPath('new', 'template-with-one-integration'), {})

      await expect(folder).toEqual({ selected: '' })
    })

    it('returns early (nested)', async () => {
      expect.assertions(1)

      const folder = await selectFolder(fixturesPath('new', 'template-with-one-integration-nested'), {})

      await expect(folder).toEqual({ selected: 'nested-here' })
    })
  })

  describe.skip('multiple integrations found', () => {
    beforeEach(() => {
      resetInquirerMock()
    })

    it('prompt to choose one from a list', async () => {
      expect.assertions(1)
      const inquirerReturn = jest.fn(() => {
        return Promise.resolve({ selected: 'ok' })
      })
      mockInquirer(inquirerReturn)

      await selectFolder(fixturesPath('new/template-with-multiple-integrations'), {})

      expect(inquirerReturn.mock.calls[0]).toMatchSnapshot()
    })

    it('does not prompt if selectedPath exist', async () => {
      expect.assertions(1)

      const { selected } = await selectFolder(fixturesPath('new/template-with-multiple-integrations'), {
        selectedPath: 'provider/2-one'
      })

      expect(selected).toEqual('provider/2-one')
    })

    it('throws an no integration found error if selectedPath does not exist', async () => {
      expect.assertions(1)

      await expect(
        selectFolder(fixturesPath('new/template-with-multiple-integrations'), { selectedPath: 'wrong-one' })
      ).rejects.toThrow('No valid integrations found under wrong-one')
    })
  })
})

describe('cloneRepository', () => {
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
