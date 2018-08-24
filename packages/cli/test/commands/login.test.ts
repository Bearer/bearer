import * as inquirer from 'inquirer'
import * as nock from 'nock'

import LoginCommand from '../../src/commands/login'

const email = 'spongebob@bearer.sh'

describe('login', () => {
  let result: Array<string>

  describe('email provided', () => {
    beforeEach(() => {
      result = []
      jest.spyOn(process.stdout, 'write').mockImplementation(val => result.push(val))

      nock('https://int.bearer.sh')
        .post('/api/v1/login', { Username: email, Password: 'ok' })
        .reply(200, {
          authorization: {
            AuthenticationResult: {
              ExpiresIn: 3600
            }
          }
        })
    })

    describe('Prompt for token', () => {
      it('logs successfully', async () => {
        jest.spyOn(inquirer, 'prompt').mockImplementation(() => Promise.resolve({ token: 'ok' }))
        await LoginCommand.run(['--email', email])
        expect(result.join()).toContain('Successfully logged in as spongebob@bearer.sh')
      })
    })

    describe('Prompt for email and token', () => {
      it('logs successfully', async () => {
        jest.spyOn(inquirer, 'prompt').mockImplementation(() => Promise.resolve({ string: email, token: 'ok' }))

        await LoginCommand.run([])
        expect(result.join()).toContain('Successfully logged in as spongebob@bearer.sh')
      })
    })
  })
})
