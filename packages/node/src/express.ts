import { NextFunction, Request, RequestHandler, Response, Router } from 'express'
import onHeaders from 'on-headers'
import bodyParser from 'body-parser'

import { CustomError } from './errors'
import Cipher from '@bearer/security'

export default (handlers: TWebhookHandlers, options: TWebhookOptions = {}) => {
  const router = Router()
  router.use(bodyParser.urlencoded({ extended: false }))
  router.use(bodyParser.json())

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
      const handler = req.bearerHandler(req)
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
  bearerHandler<T extends Request>(req: T): Promise<any>
}

type TTimedRequest = {
  startAt: number
}

export type TWebhookHandlers = Record<string, <T extends Request>(request: T) => Promise<any>>
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
    const message = JSON.stringify(req.body)
    const calculatedSha = new Cipher(token).digest(message)
    const sha = req.headers[BEARER_SHA]

    if (sha === calculatedSha) {
      next()
    } else {
      const error = new WebhookIncorrectSignature('')
      res.status(401).json({ error })
    }
  }
}
