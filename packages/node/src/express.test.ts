import express from 'express'
import request from 'supertest'

import middleware, { TWebhookHandlers } from './express'

const SUCCESS_HANDLER = 'sponge-bob-scenario-handler'
const REJECTED_HANDLER = 'patrick-is-rejecting'
const FAILING_HANDLER = 'patrick-is-failing'

const webHookHandlers: TWebhookHandlers = {
  [SUCCESS_HANDLER]: jest.fn(() => Promise.resolve(true)),
  [REJECTED_HANDLER]: jest.fn(() => Promise.reject()),
  [FAILING_HANDLER]: jest.fn(() => {
    throw new Error('ok')
    return new Promise(() => {
      throw new Error('ok')
    })
  })
}

describe('Bearer middleware', () => {
  beforeEach(() => {
    Object.keys(webHookHandlers).map(key => {
      ;(webHookHandlers[key] as any).mockClear()
    })
  })

  describe('middleware logic', () => {
    const app = express()
    app.use('/whatever/webhooks', middleware(webHookHandlers))

    app.post('/*', (_req, res) => {
      res.status(200).json({ name: 'john' })
    })

    app.get('/*', (_req, res) => {
      res.status(200).json({ name: 'john' })
    })

    describe('existing handlers', () => {
      it('calls handler', async () => {
        const response = await request(app)
          .post('/whatever/webhooks')
          .set('BEARER-SCENARIO-HANDLER', SUCCESS_HANDLER)
          .expect('Content-Type', /json/)
          .expect(200)

        expect(webHookHandlers[SUCCESS_HANDLER]).toHaveBeenCalled()
        expect(response.body).toMatchObject({ ack: 'ok' })
      })

      it('perform promise even if rejecting', async () => {
        const response = await request(app)
          .post('/whatever/webhooks')
          .set('BEARER-SCENARIO-HANDLER', REJECTED_HANDLER)
          .expect('Content-Type', /json/)
          .expect(422)

        expect(webHookHandlers[REJECTED_HANDLER]).toHaveBeenCalled()
        expect(response.body).toMatchObject({ error: { name: 'Bearer:WebhookHandlerRejection' } })
      })

      it('fails gracefully if the handler raises an error', async () => {
        const response = await request(app)
          .post('/whatever/webhooks')
          .set('BEARER-SCENARIO-HANDLER', FAILING_HANDLER)
          .expect('Content-Type', /json/)
          .expect(500)

        expect(response.body).toMatchObject({ error: { name: 'Bearer:WebhookProcessingError' } })
        expect(webHookHandlers[FAILING_HANDLER]).toHaveBeenCalled()
        expect(webHookHandlers[SUCCESS_HANDLER]).not.toHaveBeenCalled()
      })
    })

    describe('non existing handler', () => {
      it('returns unprocessable entity status', async () => {
        const response = await request(app)
          .post('/whatever/webhooks')
          .set('BEARER-SCENARIO-HANDLER', 'unknow')
          .expect('Content-Type', /json/)
          .expect(422)

        expect(response.body).toMatchObject({ error: { name: 'Bearer:WebhookMissingHandler' } })
        expect(webHookHandlers[FAILING_HANDLER]).not.toHaveBeenCalled()
        expect(webHookHandlers[SUCCESS_HANDLER]).not.toHaveBeenCalled()
      })
    })

    it('handles only the route it needs', async () => {
      const getResponse = await request(app)
        .get('/spongebobdoesnotneedwebhooks')
        .expect('Content-Type', /json/)
        .expect(200)
      const postResponse = await request(app)
        .post('/spongebobdoesnotneedwebhooks')
        .expect('Content-Type', /json/)
        .expect(200)

      expect(webHookHandlers[SUCCESS_HANDLER]).not.toHaveBeenCalled()
      expect(webHookHandlers[FAILING_HANDLER]).not.toHaveBeenCalled()

      expect(postResponse.body).toMatchObject({ name: 'john' })
      expect(getResponse.body).toMatchObject({ name: 'john' })
    })
  })

  describe('body validation', () => {
    const app = express()

    describe('secret given', () => {
      app.use('/whatever/protected', middleware(webHookHandlers, { token: '1234' }))

      it('stops unsigned requests', async () => {
        const response = await request(app)
          .post('/whatever/protected')
          .expect('Content-Type', /json/)
          .expect(401)

        expect(response.body).toMatchObject({ error: { name: 'Bearer:WebhookIncorrectSignature' } })
      })

      it('process correctly signed requests', async () => {
        await request(app)
          .post('/whatever/protected')
          .set('BEARER-SCENARIO-HANDLER', SUCCESS_HANDLER)
          .set('BEARER-SHA', '1234')
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
          .set('BEARER-SCENARIO-HANDLER', SUCCESS_HANDLER)
          .set('BEARER-SHA', 'it does not matter')
          .expect('Content-Type', /json/)
          .expect(200)
      })
    })
  })
})
