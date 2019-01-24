import { TAPIKEYAuthContext, SaveState, TSaveActionEvent, TSavePromise } from '@bearer/intents'

export default class SaveSetupIntent extends SaveState implements SaveState<State, ReturnedData, any, TAPIKEYAuthContext> {
  async action(event: TSaveActionEvent<State, Params, TAPIKEYAuthContext>): TSavePromise<State, ReturnedData> {
    return { state: event.params.setup, data: [] }
  }
}

export type Params = {
  setup: any
}

export type State = {
  apiKey: string
}


export type ReturnedData = {
}
