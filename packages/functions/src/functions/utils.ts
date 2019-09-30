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

export function fetchData(
  callback: d.TLambdaCallback,
  { data, error, statusCode }: { data?: any; error?: any; statusCode?: any }
) {
  if (data) {
    sendSuccessMessage(callback, { data, statusCode })
  } else {
    sendErrorMessage(callback, { statusCode, error: error || 'Unkown error' })
  }
}

export function eventAsActionParams(event: d.TLambdaEvent) {
  return {
    context: event.context,
    params: paramsFromEvent(event)
  }
}

function paramsFromEvent(event: d.TLambdaEvent) {
  return {
    ...event.queryStringParameters,
    ...bodyFromEvent(event)
  }
}

export const BACKEND_ONLY_ERROR = {
  code: 'UNAUTHORIZED_FUNCTION_CALL',
  message:
    // tslint:disable-next-line:max-line-length
    "This function can't be called from the frontend. If you want to call APIs from the frontend, please refer to this link for more information: https://docs.bearer.sh/integration-clients/javascript#calling-apis"
}
