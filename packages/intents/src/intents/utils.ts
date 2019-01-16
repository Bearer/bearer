import { sendErrorMessage, sendSuccessMessage } from '../lambda'
import * as d from '../declaration'

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
