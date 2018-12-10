import { NextFunction, Request, RequestHandler, Response, Router } from 'express'
import onHeaders from 'on-headers'

import { CustomError } from './errors'

export default (handlers: TWebhookHandlers, options: TWebhookOptions = {}) => {
  const router = Router()

  if (options.token) {
    router.use(verifyPayload(options.token))
  }

  router.post('/', ensureHandlerExists(handlers), (async (
    req: Request & TWithHandlerReq & TTimedRequest,
    res: Response,
    _next: NextFunction
  ) => {
    const startAt = Date.now()

    // please do not transform to arrow function
    onHeaders(res, function() {
      this.setHeader(INTENT_DURATION_HEADER, Date.now() - startAt)
    })

    try {
      const handler = req.bearerHandler()
      try {
        await handler
        res.json({ ack: 'ok' })
        // tslint:disable-next-line:no-unused
      } catch (_error) {
        const error = new WebhookHandlerRejection(req.bearerHandlerName)
        res.status(422).json({ error })
      }
      // tslint:disable-next-line:no-unused
    } catch (_e) {
      const error = new WebhookProcessingError(req.bearerHandlerName)
      res.status(500).json({ error })
    }
  }) as RequestHandler)

  return router
}

class WebhookMissingHandler extends CustomError {}
class WebhookProcessingError extends CustomError {}
class WebhookHandlerRejection extends CustomError {}
class WebhookIncorrectSignature extends CustomError {}

/**
 * Types
 */
type TWithHandlerReq = {
  bearerHandlerName: string
  bearerHandler(): Promise<any>
}

type TTimedRequest = {
  startAt: number
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
const INTENT_DURATION_HEADER = 'X-BEARER-WEBHOOK-HANDLER-DURATION'
function ensureHandlerExists(handlers: TWebhookHandlers) {
  return (req: Request & Partial<TWithHandlerReq>, res: Response, next: NextFunction) => {
    const scenarioHandler = req.headers[SCENARIO_HANDLER] as string
    if (!handlers[scenarioHandler]) {
      const error = new WebhookMissingHandler(`Scenario handler not found: ${scenarioHandler}`)
      res.status(422).json({ error })
    } else {
      req.bearerHandler = handlers[scenarioHandler]
      req.bearerHandlerName = scenarioHandler
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
      const error = new WebhookIncorrectSignature('')
      res.status(401).json({ error })
    }
  }
}
