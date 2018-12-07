import express from 'express'
import request from 'supertest'

import middleware from './middleware'

describe('Bearer middleware', () => {
  describe('middleware logic', () => {
    const app = express()
    const handler = jest.fn().mockImplementation(() => Promise.resolve(true))

    beforeEach(() => {
      handler.mockReset()
    })

    app.use('/whatever/webhooks', middleware(handler))
    app.post('/*', (_req, res) => {
      res.status(200).json({ name: 'john' })
    })
    app.get('/*', (_req, res) => {
      res.status(200).json({ name: 'john' })
    })

    it('calls handler', async () => {
      const response = await request(app)
        .post('/whatever/webhooks')
        .expect('Content-Type', /json/)
        .expect(200)
      expect(handler).toHaveBeenCalled()
      expect(response.body).toMatchObject({ ack: 'ok' })
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

      expect(handler).not.toHaveBeenCalled()

      expect(postResponse.body).toMatchObject({ name: 'john' })
      expect(getResponse.body).toMatchObject({ name: 'john' })
    })
  })

  describe('body validation', () => {
    const app = express()
    describe('secret given', () => {
      const protectedHandler = jest.fn().mockImplementation(() => Promise.resolve(true))
      app.use('/whatever/protected', middleware(protectedHandler, { token: '1234' }))

      it('stops unsigned requests', async () => {
        await request(app)
          .post('/whatever/protected')
          .expect('Content-Type', /json/)
          .expect(401)
      })

      it('process requests correctly signed', async () => {
        await request(app)
          .post('/whatever/protected')
          .set('BEARER-SHA', '1234')
          .expect('Content-Type', /json/)
          .expect(200)
        expect(protectedHandler).toHaveBeenCalled()
      })
    })

    describe('no secret given', () => {
      const unprotectedHnadler = jest.fn().mockImplementation(() => Promise.resolve(true))
      app.use('/whatever/unprotected', middleware(unprotectedHnadler))

      it('process all requests', async () => {
        await request(app)
          .post('/whatever/unprotected')
          .set('BEARER-SHA', 'it does not matter')
          .expect('Content-Type', /json/)
          .expect(200)

        expect(unprotectedHnadler).toHaveBeenCalled()
      })
    })
  })
})
