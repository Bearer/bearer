import { {{contextAuthTypeImport}}SaveState, TSaveActionEvent, TSavePromise } from '@bearer/intents'

export default class SaveSetupIntent extends SaveState implements SaveState<State, ReturnedData, any{{contextAuthTypeImplements}}> {
  async action(event: TSaveActionEvent<State, Params{{contextAuthTypeImplements}}>): TSavePromise<State, ReturnedData> {
    return { state: event.params.setup, data: [] }
  }
}

export type Params = {
  setup: any
}

export type State = {{dataType}}


export type ReturnedData = {
}
