import https from 'https'
import nock from 'nock'

import { overrideRequestMethod } from '../src/http-overrides'
import { setupFunctionIdentifiers } from '../src'

const billingLogger = jest.fn()
const userLogger = jest.fn()
Date.now = jest.fn(() => 5)

beforeAll(() => {
  setupFunctionIdentifiers({
    context: {
      clientId: 'a-client-id',
      integrationUuid: 'custom0buid'
    }
  })
  process.env._X_AMZN_TRACE_ID = 'aws_trace'
})

describe('overrideRequestMethod', () => {
  describe('https', () => {
    const url = 'https://bearer.sh'

    beforeAll(() => {
      nock(url)
        .get('/path?ok=false')
        .reply(201, { great: 'sponge bob is the king' })
      nock(url)
        .post('/postPath?postQuery=sponge')
        .reply(204, { great: 'sponge bob is the king of posts' })

      overrideRequestMethod(https, billingLogger, userLogger)
    })

    beforeEach(() => {
      billingLogger.mockClear()
      userLogger.mockClear()
    })

    describe('get', () => {
      it('logs request', async () => {
        mockDuration(10)

        nock(url)
          .get('/path?ok=false')
          .reply(201, { great: 'sponge bob is the king' })

        await new Promise(resolve => {
          https.get(`${url}/path?ok=false`, res => {
            res.on('end', resolve)
          })
        })

        expect(billingLogger).toHaveBeenCalledTimes(1)
        expect(billingLogger.mock.calls[0]).toMatchSnapshot()
        expect(userLogger).toHaveBeenCalledTimes(1)
        expect(userLogger.mock.calls[0]).toMatchSnapshot()
      })
    })

    describe('request', () => {
      it('logs request', async () => {
        mockDuration(15)

        nock(url)
          .post('/postPath?postQuery=sponge')
          .reply(201, { great: 'sponge bob is the king' })

        const data = JSON.stringify({
          todo: 'Buy the milk'
        })

        const options = {
          hostname: 'bearer.sh',
          path: '/postPath?postQuery=sponge',
          method: 'POST',
          headers: {
            UserAgent: 'Bearer'
          }
        }

        await new Promise((resolve, reject) => {
          const req = https.request(options, res => {
            res.on('end', resolve)
            res.on('error', reject)
          })
          req.write(data)
          req.end()
        })

        expect(billingLogger).toHaveBeenCalledTimes(1)
        expect(billingLogger.mock.calls[0]).toMatchSnapshot()
        expect(userLogger).toHaveBeenCalledTimes(1)
        expect(userLogger.mock.calls[0]).toMatchSnapshot()
      })
    })

    describe('special urls', () => {
      const url = 'https://logs.eu-west-3.amazonaws.com'

      beforeAll(() => {
        nock(url)
      })

      it('does not log request', async () => {
        mockDuration(5)

        nock(url)
          .post('/whatever')
          .reply(201, { great: 'sponge bob is the king' })

        const data = JSON.stringify({
          todo: 'Buy the milk'
        })

        const options = {
          hostname: 'logs.eu-west-3.amazonaws.com',
          port: 443,
          path: '/whatever',
          method: 'POST',
          headers: {
            UserAgent: 'Bearer'
          }
        }

        await new Promise((resolve, reject) => {
          const req = https.request(options, res => {
            res.on('end', resolve)
            res.on('error', reject)
          })
          req.write(data)
          req.end()
        })

        expect(billingLogger).not.toHaveBeenCalled()
        expect(userLogger).not.toHaveBeenCalled()
      })
    })
  })
})

function mockDuration(duration: number) {
  // @ts-ignore
  Date.now.mockReturnValueOnce(10)
  // @ts-ignore
  Date.now.mockReturnValueOnce(10 + duration)
}
