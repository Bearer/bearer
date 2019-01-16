import { AxiosResponse } from 'axios'

/**
 * Simplest error payload, better add code and message so that debugging is made easier
 */
export type TErrorPayload<ReturnedError> = { error: ReturnedError }
export type TDataPayload<ReturnedData> = { data: ReturnedData }

/**
 * Intent callbacks
 */
export type TFetchPayload<ReturnedData = any, ReturnedError = any> = Partial<TDataPayload<ReturnedData>> &
  Partial<TErrorPayload<ReturnedError>>

export type TRetrieveStatePayload<ReturnedData = any, ReturnedError = any> = TFetchPayload<ReturnedData, ReturnedError>

export type TSaveStatePayload<State = any, ReturnedData = any, ReturnedError = any> = {
  state: State
  data?: ReturnedData
} & Partial<TErrorPayload<ReturnedError>>

/**
 * Contexts
 */

type TBaseAuthContext<TAuthAccessContent> = { authAccess: TAuthAccessContent; [key: string]: any }
export type TOAUTH2AuthContext = TBaseAuthContext<{ accessToken: string }>
export type TNONEAuthContext = TBaseAuthContext<undefined>
export type TBASICAuthContext = TBaseAuthContext<{ username: string; password: string }>
export type TAPIKEYAuthContext = TBaseAuthContext<{ apiKey: string }>

export type TAuthContext = TNONEAuthContext | TOAUTH2AuthContext | TBASICAuthContext | TAPIKEYAuthContext

export type TBearerLambdaContext<AuthContext = TAuthContext> = AuthContext & {
  bearerBaseURL: string
}

/**
 * Intents
 */

/**
 * Save state action, let you store data into Bearer database without having to deal with database communication
 * Later, data could be automatically loaded by passing a reference ID parameter
 * terraformerId => will inject terrafomer object into context if found within Bearer database
 */
export type ISaveStateAction<AuthContext = TAuthContext, State = any, ReturnedData = any, Params = any> = (
  event: {
    context: TBearerLambdaContext<AuthContext>
    params: Params
    state: State
  }
) => Promise<TSaveStatePayload<State, ReturnedData>>

/**
 * Retrieve state action, let you retrieve data stored into Bearer database
 * Alternatively, you can retrieve data from a fetch Intent by
 */
export type IRetrieveStateAction<AuthContext = TAuthContext, State = any, ReturnedData = any, Params = any> = (
  event: {
    context: TBearerLambdaContext<AuthContext>
    params: Params
    state: State
  }
) => Promise<TFetchPayload<ReturnedData>>

/**
 * Fetch any data
 */
export type TFetchAction<AuthContext = TAuthContext, ReturnedData = any> = (
  event: {
    context: TBearerLambdaContext<AuthContext>
    params: Record<string, any>
  }
) => Promise<TFetchPayload<ReturnedData>>

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
  callback: (result: TSaveStatePayload) => void
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
  callback: (result: TRetrieveStatePayload) => void
) => void

/**
 * @deprecated since version beta5 please use IRetrieveStateAction
 */
export type TFetchDataAction = (
  context: TBearerLambdaContext,
  params: Record<string, any>,
  body: Record<string, any>,
  callback: (result: TFetchPayload) => void
) => void

/**
 * when success, state represent the data you want to store within Bearer database
 * whereras data sent to the frontend could be different
 */
export type TSaveStateCallback = (payload: TSaveStatePayload & TErrorPayload<any>) => void
