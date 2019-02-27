import * as path from 'path'

import NewCommand from '../../src/commands/new'
import { readFile as _readFile } from '../helpers/utils'

const destination = path.join(__dirname, '..', '..', '.bearer/init')
const AUTHCONFIG = 'auth.config.json'

function readFile(folder: string, filename: string): string {
  return _readFile(destination, folder, filename)
}

const auths = ['OAUTH2', 'BASIC', 'APIKEY', 'NONE', 'OAUTH1']

describe.each(Object.values(auths))('%s', auth => {
  it(`generates a scenario without any prompt and ${auth}`, async () => {
    const result: string[] = []
    jest.spyOn(process.stdout, 'write').mockImplementation(val => result.push(val))

    await NewCommand.run([
      `${auth}Scenario`,
      '-a',
      auth,
      '--skipInstall',
      '--path',
      path.join(destination, 'new', auth)
    ])

    const outPut = result.sort().join()
    expect(outPut).toMatchSnapshot()
    expect(readFile(path.join('new', auth), AUTHCONFIG)).toMatchSnapshot()
    await setTimeout(() => {}, 10)
  })
})
