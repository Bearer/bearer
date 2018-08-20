import { AxiosResponse } from 'axios'

/**
 * Intent callbacks
 */
export type TFetchDataCallback = (payload: { data: any } | { error: any }) => void

export type TRetrieveStateCallback = (payload: { data: any } | { error: any }) => void

/**
 * when success, state represent the data you want to store within Bearer database
 * whereras data sent to the frontend could be different
 */
export type TSaveStateCallback = (payload: { state: any; data?: any } | { error: any }) => void

/**
 * Intents
 */

export type TBearerLambdaContext = {
  authAccess: TAuthContext
}

export type ISaveStateIntentAction = (
  context: TBearerLambdaContext,
  _params: any,
  body: any,
  state: any,
  callback: TSaveStateCallback
) => void

export type IRetrieveStateIntentAction = (
  context: TBearerLambdaContext,
  _params: any,
  state: any,
  callback: TRetrieveStateCallback
) => void

export type TFetchDataAction = (
  context: TBearerLambdaContext,
  params: Record<string, any>,
  body: Record<string, any>,
  callback: TFetchDataCallback
) => void

/**
 * Auth definitions
 */
export type TOAUTH2AuthContext = {
  accessToken: string
  bearerBaseURL: string
  [key: string]: any
}

export type TNONEAuthContext = {
  bearerBaseURL: string
  [key: string]: any
}

export type TBASICAuthContext = {
  username: string
  password: string
  bearerBaseURL: string
  [key: string]: any
}

export type TAPIKEYAuthContext = {
  apiKey: string
  bearerBaseURL: string
  [key: string]: any
}

// Deprecated
export type TapiKeyContext = TAPIKEYAuthContext
export type Toauth2Context = TOAUTH2AuthContext
export type TbasicAuthContext = TBASICAuthContext
export type TnoAuthContext = TNONEAuthContext
export type DEPRECATEDCONTEXT = TapiKeyContext
// end Deprecated

export type TAuthContext =
  | TNONEAuthContext
  | TOAUTH2AuthContext
  | TBASICAuthContext
  | TAPIKEYAuthContext
  | DEPRECATEDCONTEXT

export type TStateData = AxiosResponse<{
  Item: any
}>

export type TLambdaEvent = {
  queryStringParameters: Record<string, any>
  context: Record<string, any> & TBearerLambdaContext
  body?: any
}

export type TLambdaCallback = (error: any | null, data: any) => void
