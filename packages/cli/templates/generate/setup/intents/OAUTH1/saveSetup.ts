import { TOAUTH1AuthContext, SaveState, TSaveActionEvent, TSavePromise } from '@bearer/intents'

export default class SaveSetupIntent extends SaveState
  implements SaveState<State, ReturnedData, any, TOAUTH1AuthContext> {
  async action(event: TSaveActionEvent<State, Params, TOAUTH1AuthContext>): TSavePromise<State, ReturnedData> {
    return { state: event.params.setup, data: [] }
  }
}

export type Params = {
  setup: any
}

export type State = {
  consumerKey: string
  consumerSecret: string
}

export type ReturnedData = {}
