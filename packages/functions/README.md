# @bearer/functions

## Usage

**Creating a FetchData function**

```ts
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
