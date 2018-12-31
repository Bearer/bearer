import { NextFunction, Request, RequestHandler, Response, Router } from 'express'
import onHeaders from 'on-headers'
import bodyParser from 'body-parser'
import Cipher from '@bearer/security'
import debug from '@bearer/logger'

const logger = debug('webhooks')
import { CustomError } from './errors'

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
    next: NextFunction
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
      } catch (error) {
        logger(error)
        res.status(422)
        next(new WebhookHandlerRejection(req.bearerHandlerName))
      }
    } catch (error) {
      logger(error)
      res.status(500)
      next(new WebhookProcessingError(req.bearerHandlerName))
    }
  }) as RequestHandler)

  // Error handler
  router.use((error: CustomError, req: Request, res: Response, next: NextFunction) => {
    logger(error)
    res.json({ error: { name: error.name, message: error.message } })
  })

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
      res.status(422)
      next(new WebhookMissingHandler(`Scenario handler not found: ${scenarioHandler}`))
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
      res.status(401)
      next(new WebhookIncorrectSignature('Incorrect signature, please make ure you provided the correct token'))
    }
  }
}
