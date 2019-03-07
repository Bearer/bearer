import { TBASICAuthContext, SaveState, TSaveActionEvent, TSavePromise } from '@bearer/functions'

export default class SaveSetupFunction extends SaveState implements SaveState<State, ReturnedData, any, TBASICAuthContext> {
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
