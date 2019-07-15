import express, { Request } from 'express'
import { IncomingMessage } from 'http'
import request from 'supertest'
import Cipher from '@bearer/security'

import middleware, { TWebhookHandlers } from '.'

const SUCCESS_HANDLER = 'sponge-bob-integration-handler'
const REJECTED_HANDLER = 'patrick-is-rejecting'
const FAILING_HANDLER = 'patrick-is-failing'

const webHookHandlers = {
  [SUCCESS_HANDLER]: jest.fn((req: Request & { title: string }) => Promise.resolve(req.title)),
  [REJECTED_HANDLER]: jest.fn((req: Request) => Promise.reject()),
  [FAILING_HANDLER]: jest.fn((req: Request) => {
    throw new Error('ok')
  })
} as TWebhookHandlers

describe('Bearer middleware', () => {
  beforeEach(() => {
    Object.keys(webHookHandlers).map(key => {
      // @ts-ignore
      webHookHandlers[key].mockClear()
    })
  })

  describe('middleware logic', () => {
    const app = express()
    app.use((req: Request & { title?: string }, _res, next) => {
      req.title = 'Sponge bob is the king'
      next()
    })
    app.use('/whatever/webhooks', middleware(webHookHandlers))

    app.post('/*', (_req, res) => {
      res.status(200).json({ name: 'john' })
    })

    app.get('/*', (_req, res) => {
      res.status(200).json({ name: 'john' })
    })

    describe('existing handlers', () => {
      it('calls handler with req object', async () => {
        await request(app)
          .post('/whatever/webhooks')
          .set('BEARER-INTEGRATION-HANDLER', SUCCESS_HANDLER)

        const mock = webHookHandlers[SUCCESS_HANDLER] as jest.Mock
        // request is passed down
        expect(mock.mock.calls[0][0]).toBeInstanceOf(IncomingMessage)
        expect(mock.mock.calls[0][0].title).toBe('Sponge bob is the king')
      })

      it('with alias - calls handler with req object', async () => {
        await request(app)
          .post('/whatever/webhooks')
          .set('BEARER-INTEGRATION-HANDLER', 'whatever')
          .set('BEARER-INTEGRATION-ALIAS', SUCCESS_HANDLER)

        const mock = webHookHandlers[SUCCESS_HANDLER] as jest.Mock
        // request is passed down
        expect(mock.mock.calls[0][0]).toBeInstanceOf(IncomingMessage)
        expect(mock.mock.calls[0][0].title).toBe('Sponge bob is the king')
      })

      it('returns correct success payload', async () => {
        const response = await request(app)
          .post('/whatever/webhooks')
          .set('BEARER-INTEGRATION-HANDLER', SUCCESS_HANDLER)
          .expect('Content-Type', /json/)
          .expect(200)

        expect(webHookHandlers[SUCCESS_HANDLER]).toHaveBeenCalled()
        expect(response.body).toMatchObject({ ack: 'ok' })
      })

      it('perform promise even if rejecting', async () => {
        const response = await request(app)
          .post('/whatever/webhooks')
          .set('BEARER-INTEGRATION-HANDLER', REJECTED_HANDLER)
          .expect('Content-Type', /json/)
          .expect(422, {
            error: { name: 'Bearer:WebhookHandlerRejection', message: 'patrick-is-rejecting' }
          })

        expect(webHookHandlers[REJECTED_HANDLER]).toHaveBeenCalled()
      })

      it('fails gracefully if the handler raises an error', async () => {
        const response = await request(app)
          .post('/whatever/webhooks')
          .set('BEARER-INTEGRATION-HANDLER', FAILING_HANDLER)
          .expect('Content-Type', /json/)
          .expect(500, {
            error: { name: 'Bearer:WebhookProcessingError', message: 'patrick-is-failing' }
          })

        expect(webHookHandlers[FAILING_HANDLER]).toHaveBeenCalled()
        expect(webHookHandlers[SUCCESS_HANDLER]).not.toHaveBeenCalled()
      })
    })

    describe('non existing handler', () => {
      it('returns unprocessable entity status', async () => {
        const response = await request(app)
          .post('/whatever/webhooks')
          .set('BEARER-INTEGRATION-HANDLER', 'unknow')
          .expect('Content-Type', /json/)
          .expect(422, {
            error: { name: 'Bearer:WebhookMissingHandler', message: 'Integration handler not found: unknow' }
          })

        expect(webHookHandlers[FAILING_HANDLER]).not.toHaveBeenCalled()
        expect(webHookHandlers[SUCCESS_HANDLER]).not.toHaveBeenCalled()
      })
    })

    it('handles only the route it needs', async () => {
      const getResponse = await request(app)
        .get('/spongebobdoesnotneedwebhooks')
        .expect('Content-Type', /json/)
        .expect(200, { name: 'john' })
      const postResponse = await request(app)
        .post('/spongebobdoesnotneedwebhooks')
        .expect('Content-Type', /json/)
        .expect(200, { name: 'john' })

      expect(webHookHandlers[SUCCESS_HANDLER]).not.toHaveBeenCalled()
      expect(webHookHandlers[FAILING_HANDLER]).not.toHaveBeenCalled()
    })
  })

  describe('body validation', () => {
    const app = express()
    const body = { key1: 'value1' }
    const masterKey = '1234'
    const cipher = new Cipher(masterKey)

    describe('secret given', () => {
      app.use('/whatever/protected', middleware(webHookHandlers, { token: masterKey }))

      it('stops unsigned requests', async () => {
        await request(app)
          .post('/whatever/protected')
          .expect('Content-Type', /json/)
          .expect(401, {
            error: {
              name: 'Bearer:WebhookIncorrectSignature',
              message: 'Incorrect signature, please make sure you provided the correct token'
            }
          })
      })

      it('stops bad signed requests', async () => {
        await request(app)
          .post('/whatever/protected')
          .set('BEARER-SHA', 'bad-signature')
          .expect('Content-Type', /json/)
          .expect(401, {
            error: {
              name: 'Bearer:WebhookIncorrectSignature',
              message: 'Incorrect signature, please make sure you provided the correct token'
            }
          })
      })

      it('process correctly signed requests', async () => {
        const signature = cipher.digest(JSON.stringify(body))
        await request(app)
          .post('/whatever/protected')
          .set('BEARER-INTEGRATION-HANDLER', SUCCESS_HANDLER)
          .set('BEARER-SHA', signature)
          .send(body)
          .expect('Content-Type', /json/)
          .expect(res => {
            return res.header['X-BEARER-WEBHOOK-HANDLER-DURATION'] > 0.0
          })
          .expect(200)
      })
    })

    describe('no secret given', () => {
      app.use('/whatever/unprotected', middleware(webHookHandlers))

      it('process all requests', async () => {
        await request(app)
          .post('/whatever/unprotected')
          .set('BEARER-INTEGRATION-HANDLER', SUCCESS_HANDLER)
          .set('BEARER-SHA', 'it does not matter')
          .expect('Content-Type', /json/)
          .expect(200)
      })
    })
  })
})
