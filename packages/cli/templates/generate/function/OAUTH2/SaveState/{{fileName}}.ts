import { TOAUTH2AuthContext, SaveState, TSaveActionEvent, TSavePromise } from '@bearer/functions'
// Uncomment this line if you need to use Client
// import Client from './client'

export default class {{functionClassName}}Function extends SaveState implements SaveState<State, ReturnedData, any, TOAUTH2AuthContext> {
  async action(event: TSaveActionEvent<State, Params, TOAUTH2AuthContext>): TSavePromise<State, ReturnedData> {
    // const token = event.context.authAccess.accessToken
    // Put your logic here
    return { state: [], data: [] }
  }

  // Uncomment the line above if you don't want your function to be called from the frontend
  // static backendOnly = true
  
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

export type ReturnedData = {
  // foo: string[]
}
