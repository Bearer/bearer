import * as d from '../declaration'
import { bodyFromEvent, fetchData } from './utils'

export class FetchData {
  static intent(action: d.TFetchDataAction) {
    return (event: d.TLambdaEvent, _context, lambdaCallback: d.TLambdaCallback) => {
      action(event.context, event.queryStringParameters, bodyFromEvent(event), result => {
        fetchData(lambdaCallback, result)
      })
    }
  }
}
