# @bearer/functions

## Usage

**Creating a FetchData function**

```
import { FetchData, TOAUTH2AuthContext } from '@bearer/functions'
import Client from './client'

export default class FunctionName extends FetchData implements FetchData<ReturnedData, any, TOAUTH2AuthContext> {
  // Uncomment the line above if you don't want your function to be called from the frontend
  // static backendOnly = true

  async action(event: TFetchActionEvent<Params, TOAUTH2AuthContext>): TFetchPromise<ReturnedData> {
    // const token = event.context.authAccess.accessToken
    // Put your logic here
    return { data: [] }
  }
}

export type Params = {
  // name: string
}

export type ReturnedData = {
  // foo: string[]
}
```

**Creating a SaveData function**

```
import { TOAUTH2AuthContext, SaveState, TSaveActionEvent, TSavePromise } from '@bearer/functions'
// Uncomment this line if you need to use Client
// import Client from './client'

export default class SaveStateFunctionName extends SaveState implements SaveState<State, ReturnedData, any, TOAUTH2AuthContext> {
  // Uncomment the line above if you don't want your function to be called from the frontend
  // static backendOnly = true

  async action(event: TSaveActionEvent<State, Params, TOAUTH2AuthContext>): TSavePromise<State, ReturnedData> {
    // const token = event.context.authAccess.accessToken
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

export type ReturnedData = {
  // foo: string[]
}

```
