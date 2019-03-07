import { TAPIKEYAuthContext, SaveState, TSaveActionEvent, TSavePromise } from '@bearer/functions'

export default class SaveSetupFunction extends SaveState implements SaveState<State, ReturnedData, any, TAPIKEYAuthContext> {
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
