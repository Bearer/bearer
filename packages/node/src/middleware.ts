import { NextFunction, Request, RequestHandler, Response, Router } from 'express'

export default (handlers: TWebhookHandlers, options: TWebhookOptions = {}) => {
  const router = Router()

  if (options.token) {
    router.use(verifyPayload(options.token))
  }

  router.post('/', ensureHandlerExists(handlers), (async (
    req: Request & TWithHandlerReq,
    res: Response,
    _next: NextFunction
  ) => {
    try {
      const handler = req.bearerHandler()
      try {
        await handler
        res.json({ ack: 'ok' })
        // tslint:disable-next-line:no-unused
      } catch (_error) {
        res.status(422).json({
          message: 'Rejecting incoming webhook'
        })
      }
      // tslint:disable-next-line:no-unused
    } catch (_e) {
      res.status(500).json({ ack: 'error' })
    }
  }) as RequestHandler)

  return router
}

class MissingWebhookHandler extends Error {
  constructor(private readonly scenarioHandlerName: string) {
    super()
  }
}

/**
 * Types
 */
type TWithHandlerReq = {
  bearerHandler(): Promise<any>
}

export type TWebhookHandlers = Record<string, () => Promise<any>>
export type TWebhookOptions = {
  token?: string
}

/**
 * Middlwares
 */
const SCENARIO_HANDLER = 'bearer-scenario-handler'
const BEARER_SHA = 'bearer-sha'

function ensureHandlerExists(handlers: TWebhookHandlers) {
  return (req: Request & Partial<TWithHandlerReq>, res: Response, next: NextFunction) => {
    const scenarioHandler = req.headers[SCENARIO_HANDLER] as string
    if (!handlers[scenarioHandler]) {
      const error = new MissingWebhookHandler(scenarioHandler)
      res.status(422).json({
        message: error.toString()
      })
    } else {
      req.bearerHandler = handlers[scenarioHandler]
      next()
    }
  }
}

function verifyPayload(token: string) {
  return (req: Request, res: Response, next: NextFunction) => {
    const sha = req.headers[BEARER_SHA]
    if (sha === token) {
      next()
    } else {
      res.statusCode = 401
      res.json({ error: 'Incorrect signature' })
    }
  }
}
