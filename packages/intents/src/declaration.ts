import { AxiosResponse } from 'axios'

export type TIntentSuccessCallback = {
  ({ data }: { data: any }): void
}

export type TRetrieveStateSuccessCallback = {
  ({ state }: { state: any }): void
}

export type TErrorCallback = {
  ({ error }: { error: any }): void
}

export type TSaveStateSuccessCallback = {
  ({ state, data }: { state: any; data: any }): void
}

export type TIntentCallback = TIntentSuccessCallback | TErrorCallback
export type TFetchDataCallback = TIntentCallback

export type TRetrieveStateCallback = TRetrieveStateSuccessCallback | TErrorCallback

export type TSaveStateCallback = TSaveStateSuccessCallback | TErrorCallback

/**
 * Intents
 */

export type TBearerLambdaContext = {
  authAccess: TAuthContext
}

export type ISaveStateIntentAction = {
  (context: TBearerLambdaContext, _params: any, body: any, state: any, callback: TSaveStateCallback): void
}

export type IRetrieveStateIntentAction = {
  (context: TBearerLambdaContext, _params: any, state: any, callback: TRetrieveStateCallback): void
}

export type TFetchDataAction = {
  (
    context: TBearerLambdaContext,
    params: Record<string, any>,
    body: Record<string, any>,
    callback: TIntentCallback
  ): void
}

/**
 * Auth definitions
 */
export type Toauth2Context = {
  accessToken: string
  bearerBaseURL: string
  [key: string]: any
}

export type TnoAuthContext = {
  bearerBaseURL: string
  [key: string]: any
}

export type TbasicAuthContext = {
  username: string
  password: string
  bearerBaseURL: string
  [key: string]: any
}

export type TapiKeyContext = {
  apiKey: string
  bearerBaseURL: string
  [key: string]: any
}

export type TAuthContext = TnoAuthContext | Toauth2Context | TbasicAuthContext | TapiKeyContext

export type TStateData = AxiosResponse<{
  Item: any
}>

export type TLambdaEvent = {
  queryStringParameters: Record<string, any>
  context: Record<string, any> & TBearerLambdaContext
  body?: string
}

export type TLambdaCallback = {
  (error: any | null, data: any): void
}
