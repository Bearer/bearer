import { TBASICAuthContext, TSaveStatePayload, SaveState, TSaveActionEvent } from '@bearer/intents'
// Uncomment this line if you need to use Client
// import Client from './client'

export default class {{intentClassName}}Intent {
  intentType = SaveState

  async action(event: TSaveActionEvent<TBASICAuthContext, State, Params>): Promise<ReturnedData> {
    // const { username, password } = event.context.authAccess
    // Put your logic here
    return { state: [], data: [] }
  }
}

/**
 * Typing
 */
export type Params = {
  // name: string
}

export type State = {
  // foo: string[]
}

export type ReturnedData = TSaveStatePayload<State, {
  // foo: string[]
}>
