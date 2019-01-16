import { AxiosResponse } from 'axios'

/**
 * Simplest error payload, better add code and message so that debugging is made easier
 */
export type TErrorPayload = { error: any }
export type TDataPayload<ReturnedData = any> = { data: ReturnedData }

/**
 * Intent callbacks
 */
export type TFetchPayload<ReturnedData = any> = (payload: TDataPayload<ReturnedData> | TErrorPayload) => void
export type TRetrieveStatePayload<ReturnedData = any> = (payload: TDataPayload<ReturnedData> | TErrorPayload) => void

export type TSaveStatePayload<State = any, ReturnedData = any> = {
  state: State
  data?: ReturnedData
}

/**
 * Intents
 */
export type TBearerLambdaContext<AuthContext = TAuthContext> = AuthContext & {
  bearerBaseURL: string
}

/**
 * Save state action, let you store data into Bearer database without having to deal with database communication
 * Later, data could be automatically loaded by passing a reference ID parameter
 * terraformerId => will inject terrafomer object into context if found within Bearer database
 */
export type ISaveStateAction<AuthContext = TAuthContext, State = any> = (
  context: TBearerLambdaContext<AuthContext>,
  params: any,
  body: any,
  state: State
) => Promise<TSaveStatePayload>

/**
 * Retrieve state action, let you retrieve data stored into Bearer database
 * Alternatively, you can retrieve data from a fetch Intent by
 */
export type IRetrieveStateAction<AuthContext = TAuthContext, State = any> = (
  context: TBearerLambdaContext<AuthContext>,
  params: any,
  state: State
) => Promise<TFetchPayload>

/**
 * Fetch any data
 */
export type TFetchAction<AuthContext = TAuthContext, ReturnedData = any> = (
  context: TBearerLambdaContext<AuthContext>,
  params: Record<string, any>,
  body: Record<string, any>
) => Promise<TFetchPayload<ReturnedData>>

/**
 * Auth definitions
 */
type TBaseAuthContext<TAuthAccessContent> = {
  authAccess: TAuthAccessContent
  [key: string]: any
}

export type TOAUTH2AuthContext = TBaseAuthContext<{ accessToken: string }>

export type TNONEAuthContext = TBaseAuthContext<undefined>

export type TBASICAuthContext = TBaseAuthContext<{
  username: string
  password: string
}>

export type TAPIKEYAuthContext = TBaseAuthContext<{ apiKey: string }>

export type TAuthContext = TNONEAuthContext | TOAUTH2AuthContext | TBASICAuthContext | TAPIKEYAuthContext

export type TStateData = AxiosResponse<{
  Item: any
}>

export type TLambdaEvent<T = TAuthContext> = {
  queryStringParameters: Record<string, any>
  context: Record<string, any> & TBearerLambdaContext<T>
  body?: any
}

export type TLambdaCallback = (error: any | null, data: any) => void

/**
 * Deprecated
 */

/**
 * @deprecated since version beta5 please use ISaveStateAction
 *
 * Save state action, let you store data into Bearer database without having to deal with database communication
 * Later, data could be automatically loaded by passing a reference ID parameter
 * terraformerId => will inject terrafomer object into context if found within Bearer database
 */
export type ISaveStateIntentAction = (
  context: TBearerLambdaContext,
  _params: any,
  body: any,
  state: any,
  callback: TSaveStateCallback
) => void

/**
 * @deprecated since version beta5 please use IRetrieveStateAction
 * Retrieve state action, let you retrieve data stored into Bearer database
 * Alternatively, you can retrieve data from a fetch Intent by
 */
export type IRetrieveStateIntentAction = (
  context: TBearerLambdaContext,
  _params: any,
  state: any,
  callback: TRetrieveStatePayload
) => void

/**
 * @deprecated since version beta5 please use IRetrieveStateAction
 */
export type TFetchDataAction = (
  context: TBearerLambdaContext,
  params: Record<string, any>,
  body: Record<string, any>,
  callback: TFetchPayload
) => void

/**
 * when success, state represent the data you want to store within Bearer database
 * whereras data sent to the frontend could be different
 */
export type TSaveStateCallback = (payload: TSaveStatePayload | TErrorPayload) => void
