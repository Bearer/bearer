import { NextFunction, Request, Response, Router } from 'express'

export type TWebhookHandler = () => Promise<any>
export type TWebhookOptions = {
  token?: string
}

export default (handler: TWebhookHandler, options: TWebhookOptions = {}) => {
  const router = Router()

  if (options.token) {
    router.use(verifyPayload(options.token))
  }

  router.post('/', async (_req: Request, res: Response, _next: NextFunction) => {
    try {
      await handler()
      res.json({
        ack: 'ok'
      })
      // tslint:disable-next-line:no-unused
    } catch (_e) {
      res.statusCode = 500
      res.json({
        ack: 'error'
      })
    }
  })
  return router
}

function verifyPayload(token: string) {
  return (req: Request, res: Response, next: NextFunction) => {
    if (!token) {
      next()
    } else {
      const sha = req.headers['bearer-sha']
      if (sha === token) {
        next()
      } else {
        res.statusCode = 401
        res.json({ error: 'Incorrect signature' })
      }
    }
  }
}
