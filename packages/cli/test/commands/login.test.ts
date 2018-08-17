import { test } from '@oclif/test'
import { expect } from 'fancy-test'
import * as inquirer from 'inquirer'
import * as sinon from 'sinon'

const email = 'spongebob@bearer.sh'

describe('login', () => {
  afterEach(() => (inquirer.prompt as any).restore())
  describe('email provided', () => {
    beforeEach(() => {
      sinon
        .stub(inquirer, 'prompt')
        .onCall(0)
        .resolves({ token: 'ok' })
    })
    describe('Prompt for token', () => {
      // it('works', done => {
      test
        .nock('https://int.bearer.sh', api =>
          api.post('/api/v1/login', { Username: email, Password: 'ok' }).reply(200, {
            authorization: {
              AuthenticationResult: {
                ExpiresIn: 3600
              }
            }
          })
        )
        .stdout()
        .command(['login', '--email', email])
        .it('Display success message', ctx => {
          expect(ctx.stdout).to.contain('Successfully logged in as spongebob@bearer.sh')
        })
    })
  })

  describe('email not provided', () => {
    beforeEach(() => {
      const stub = sinon.stub(inquirer, 'prompt')
      stub.resolves(Promise.resolve({ email, token: 'ok' }))
    })
    describe('Prompt for email and token', () => {
      test
        .nock('https://int.bearer.sh', api =>
          api.post('/api/v1/login', { Username: email, Password: 'ok' }).reply(200, {
            authorization: {
              AuthenticationResult: {
                ExpiresIn: 3600
              }
            }
          })
        )
        .stdout()
        .command(['login'])
        .it('Display success message', ctx => {
          expect(ctx.stdout).to.contain('Successfully logged in as spongebob@bearer.sh')
        })
    })
  })
})
