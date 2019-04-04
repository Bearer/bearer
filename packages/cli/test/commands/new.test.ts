import * as path from 'path'

import NewCommand from '../../src/commands/new'
import { readFile as _readFile } from '../helpers/utils'

const destination = path.join(__dirname, '..', '..', '.bearer/init')
const destinationWithViews = path.join(__dirname, '..', '..', '.bearer/initWithView')

const AUTHCONFIG = 'auth.config.json'
const PACKAGE_JSON = 'package.json'

const auths = ['OAUTH2', 'BASIC', 'APIKEY', 'NONE', 'OAUTH1']

describe.each(Object.values(auths))('without views %s', auth => {
  it(`generates an integration without any prompt and ${auth}`, async () => {
    const result: string[] = []
    jest.spyOn(process.stdout, 'write').mockImplementation(val => {
      result.push(val)
      return true
    })

    const out = path.join(destination, 'new', auth)

    await NewCommand.run([`${auth}Integration`, '-a', auth, '--skipInstall', '--path', out])

    const outPut = result.sort().join()
    expect(outPut).toMatchSnapshot()
    expect(_readFile(out, AUTHCONFIG)).toMatchSnapshot()

    expect(_readFile(out, PACKAGE_JSON)).toMatchSnapshot()
    await setTimeout(() => {}, 10)
  })
})

describe.each(Object.values(auths))('with views %s', auth => {
  it(`generates an integration without any prompt and ${auth}`, async () => {
    const result: string[] = []
    jest.spyOn(process.stdout, 'write').mockImplementation(val => {
      result.push(val)
      return true
    })
    const out = path.join(destinationWithViews, 'new', auth)

    await NewCommand.run([`${auth}Integration`, '-a', auth, '--skipInstall', '--withViews', '--path', out])

    const outPut = result.sort().join()

    expect(outPut).toMatchSnapshot()
    expect(_readFile(out, AUTHCONFIG)).toMatchSnapshot()

    expect(_readFile(out, PACKAGE_JSON)).toMatchSnapshot()
  })
})
