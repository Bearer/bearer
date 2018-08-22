import { expect, test } from '@oclif/test'
import * as fs from 'fs-extra'
import * as path from 'path'

const destination = path.join(__dirname, '..', '.bearer/init')

function emptyInitFolders() {
  if (fs.existsSync(destination)) {
    fs.emptyDirSync(destination)
  } else {
    fs.mkdirpSync(destination)
  }
}

function readFile(folder: string, filename: string): string {
  return fs.readFileSync(path.join(destination, folder, filename), { encoding: 'utf8' })
}

const AUTHCONFIG = 'auth.config.json'

describe('new command', () => {
  emptyInitFolders()
  test
    .stdout()
    .command(['new', 'Oauth2Scenario', '-a', 'OAUTH2', '--skipInstall', '--path', path.join(destination, 'oauth2')])
    .it('generates a scenario without any prompt and OAuth2', ctx => {
      console.log('out', ctx.stdout)
      expect(ctx.stdout).to.contain(`create: test/.bearer/init/oauth2/LICENSE`)
      expect(ctx.stdout).to.contain(`create: test/.bearer/init/oauth2/intents/tsconfig.json`)
      expect(ctx.stdout).to.contain(`create: test/.bearer/init/oauth2/auth.config.json`)
      expect(ctx.stdout).to.contain(`Create setup files`)
      expect(ctx.stdout).to.contain(`What's next?`)

      expect(readFile('oauth2', AUTHCONFIG)).to.contain('authType": "OAUTH2",')
    })

  test
    .stdout()
    .command(['new', 'NoAuthScenario', '-a', 'NONE', '--skipInstall', '--path', path.join(destination, 'none')])
    .it('generates a scenario without any prompt  and NoAuth', ctx => {
      console.log('out', ctx.stdout)
      expect(ctx.stdout).to.contain(`create: test/.bearer/init/none/LICENSE`)
      expect(ctx.stdout).to.contain(`create: test/.bearer/init/none/intents/tsconfig.json`)
      expect(ctx.stdout).to.contain(`create: test/.bearer/init/none/auth.config.json`)
      expect(ctx.stdout).to.contain(`Create setup files`)
      expect(ctx.stdout).to.contain(`What's next?`)

      expect(readFile('none', AUTHCONFIG)).to.contain('authType": "NONE"')
    })

  test
    .stdout()
    .command(['new', 'NoAuthScenario', '-a', 'APIKEY', '--skipInstall', '--path', path.join(destination, 'apikey')])
    .it('generates a scenario without any prompt and ApiKey', ctx => {
      console.log('out', ctx.stdout)
      expect(ctx.stdout).to.contain(`create: test/.bearer/init/apikey/LICENSE`)
      expect(ctx.stdout).to.contain(`create: test/.bearer/init/apikey/intents/tsconfig.json`)
      expect(ctx.stdout).to.contain(`create: test/.bearer/init/apikey/auth.config.json`)
      expect(ctx.stdout).to.contain(`Create setup files`)
      expect(ctx.stdout).to.contain(`What's next?`)

      expect(readFile('apikey', AUTHCONFIG)).to.contain('authType": "APIKEY"')
    })

  test
    .stdout()
    .command(['new', 'NoAuthScenario', '-a', 'BASIC', '--skipInstall', '--path', path.join(destination, 'basic')])
    .it('generates a scenario without any prompt and Basic Auth', ctx => {
      console.log('out', ctx.stdout)
      expect(ctx.stdout).to.contain(`create: test/.bearer/init/basic/LICENSE`)
      expect(ctx.stdout).to.contain(`create: test/.bearer/init/basic/intents/tsconfig.json`)
      expect(ctx.stdout).to.contain(`create: test/.bearer/init/basic/auth.config.json`)
      expect(ctx.stdout).to.contain(`Create setup files`)
      expect(ctx.stdout).to.contain(`What's next?`)

      expect(readFile('basic', AUTHCONFIG)).to.contain('authType": "BASIC"')
    })
})
