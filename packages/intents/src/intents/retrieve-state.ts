import * as d from '../declaration'
import { DBClient as CLIENT } from '../db-client'

export class RetrieveState {
  static intent(action: d.IRetrieveStateIntentAction) {
    // tslint:disable-next-line:variable-name
    const DBClient = CLIENT.instance

    return (event: d.TLambdaEvent, _context: any, lambdaCallback: d.TLambdaCallback): void => {
      const { referenceId } = event.queryStringParameters
      try {
        DBClient(event.context.signature)
          .getData(referenceId)
          .then(state => {
            if (state) {
              action(
                event.context,
                event.queryStringParameters,
                state.Item,
                (recoveredState: { data: any }): void => {
                  lambdaCallback(null, { meta: { referenceId: state.Item.referenceId }, data: recoveredState.data })
                }
              )
            } else {
              lambdaCallback(null, { statusCode: 404, body: JSON.stringify({ referenceId, error: 'No data found' }) })
            }
          })
      } catch (error) {
        lambdaCallback(error.toString(), { error: error.toString() })
      }
    }
  }
}
