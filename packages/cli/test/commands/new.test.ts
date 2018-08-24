import * as fs from 'fs-extra'
import * as path from 'path'

import NewCommand from '../../src/commands/new'
import { readFile as _readFile } from '../helpers/utils'

const destination = path.join(__dirname, '..', '.bearer/init')
const AUTHCONFIG = 'auth.config.json'

function emptyInitFolders() {
  if (fs.existsSync(destination)) {
    fs.emptyDirSync(destination)
  } else {
    fs.mkdirpSync(destination)
  }
}

function readFile(folder: string, filename: string): string {
  return _readFile(destination, folder, filename)
}

describe('new command', () => {
  let result: Array<string>

  beforeEach(() => {
    result = []
    jest.spyOn(process.stdout, 'write').mockImplementation(val => result.push(val))
    emptyInitFolders()
  })
  const auths = ['OAUTH2', 'BASIC', 'APIKEY', 'NONE']
  auths.map(auth => {
    it(`generates a scenario without any prompt and ${auth}`, async () => {
      await NewCommand.run([`${auth}Scenario`, '-a', auth, '--skipInstall', '--path', path.join(destination, auth)])
      const outPut = result.sort().join()
      expect(outPut).toMatchSnapshot()
      expect(readFile(auth, AUTHCONFIG)).toMatchSnapshot()
    })
  })
})
