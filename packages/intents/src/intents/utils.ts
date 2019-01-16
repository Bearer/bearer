import { sendErrorMessage, sendSuccessMessage } from '../lambda'
import * as d from '../declaration'

/**
 * We should remove this at some point, we had a difference between IS and  the local
 * web server.
 * @param event
 */
export function bodyFromEvent(event: d.TLambdaEvent): any {
  const { body } = event
  if (!body) {
    return {}
  }
  if (typeof body === 'string') {
    return JSON.parse(body)
  }
  return body
}

export function fetchData(callback: d.TLambdaCallback, { data, error }: { data?: any; error?: any }) {
  if (data) {
    sendSuccessMessage(callback, { data })
  } else {
    sendErrorMessage(callback, { error: error || 'Unkown error' })
  }
}

export function paramsFromEvent(event: d.TLambdaEvent) {
  return {
    ...event.queryStringParameters,
    ...bodyFromEvent(event)
  }
}
