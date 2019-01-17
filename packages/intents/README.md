# @bearer/intents

## Usage

**Creating a FetchData intent**

```
import { FetchData, TOAUTH2AuthContext } from '@bearer/intents'
import Client from './client'

export default class ListPullRequestsIntent {
  static intentName: string = 'listPullRequests'
  static intentType: any = FetchData

  static async action({ context, params }: { context: TOAUTH2AuthContext; params: any }) {
    try {
      const { data } = await Client(context.authAccess.accessToken).get(`/repos/${params.fullName}/pulls`, {
        params: { per_page: 10, ...params }
      })
      return { data }
    } catch (error) {
      return { error: error.toString() }
    }
  }
}

```

**Creating a SaveData intent**

```
import { SaveState, TOAUTH2AuthContext } from '@bearer/intents'
import Client from './client'

export type TState = {
  pullRequests: any
}

export type TParams = {
  fullName: string,
  page?: number
}

export default class ListPullRequestsIntent {
  static intentName: string = 'listPullRequests'
  static intentType: any = SaveState

  static async action({ context, params, state }: { context: TOAUTH2AuthContext; params: TParams, state: TState }) {
    try {
      const { data } = await Client(context.authAccess.accessToken).get(`/repos/${params.fullName}/pulls`, {
        params: { per_page: 10, ...params }
      })
      return { data, state: { pullRequests: data } }
    } catch (error) {
      return { error: error.toString() }
    }
  }
}

```
