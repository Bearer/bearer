# @bearer/functions

[![Version](https://img.shields.io/npm/v/@bearer/functions.svg)](https://npmjs.org/package/@bearer/functions)
![npm bundle size (scoped)](https://img.shields.io/bundlephobia/minzip/@bearer/functions.svg)
[![Downloads/week](https://img.shields.io/npm/dw/@bearer/functions.svg)](https://npmjs.org/package/@bearer/functions)
[![License](https://img.shields.io/npm/l/@bearer/functions.svg)](https://github.com/Bearer/bearer/packages/cli/blob/master/package.json)

## Usage

**Creating a FetchData function**

```ts
import { FetchData, TOAUTH2AuthContext } from '@bearer/functions'
import Client from './client'

export default class FunctionName extends FetchData implements FetchData<ReturnedData, any, TOAUTH2AuthContext> {
  // Uncomment the line above if you don't want your function to be called from the frontend
  // static serverSideRestricted = true

  async action(event: TFetchActionEvent<Params, TOAUTH2AuthContext>): TFetchPromise<ReturnedData> {
    // const token = event.context.auth.accessToken
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
