import { TBASICAuthContext, SaveState, TSaveActionEvent, TSavePromise } from '@bearer/intents'

export default class SaveSetupIntent extends SaveState implements SaveState<State, ReturnedData, any, TBASICAuthContext> {
  async action(event: TSaveActionEvent<State, Params, TBASICAuthContext>): TSavePromise<State, ReturnedData> {
    return { state: event.params.setup, data: [] }
  }
}

export type Params = {
  setup: any
}

export type State = {
  username: string
  password: string
}

export type ReturnedData = {}
