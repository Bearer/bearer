import * as nock from 'nock'
import LoginCommand from '../../src/commands/login'
// jest.mock('open')

describe.skip('login', () => {
  let result: string[]

  describe('email provided', () => {
    beforeEach(() => {
      result = []
      // jest.spyOn(process.stdout, 'write').mockImplementation(val => result.push(val))

      nock('https://int.bearer.sh')
        .post('/api/v1/login')
        .reply(200, {
          authorization: {
            AuthenticationResult: {
              ExpiresIn: 3600
            }
          }
        })
    })

    describe('O', () => {
      it('logs successfully', async () => {
        await LoginCommand.run([])

        expect(result.join()).toContain('Successfully logged in as spongebob@bearer.sh')
      })
    })
  })
})
