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
  body?: any
}

export type TLambdaCallback = (error: any | null, data: any) => void
