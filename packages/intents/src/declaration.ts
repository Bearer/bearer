import { AxiosResponse } from 'axios'

/**
 * Simplest error payload, better add code and message so that debugging is made easier
 */
export type TErrorPayload<ReturnedError> = { error: ReturnedError }
export type TDataPayload<ReturnedData> = { data: ReturnedData }

export type TFetchPayload<ReturnedData = any, ReturnedError = any> = Partial<TDataPayload<ReturnedData>> &
  Partial<TErrorPayload<ReturnedError>>

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
export type TSaveStateAction<AuthContext = TAuthContext, State = any, ReturnedData = any, Params = any> = (
  event: TSaveActionEvent<AuthContext, State, Params>
) => Promise<TSaveStatePayload<State, ReturnedData>>

export type TSaveActionEvent<AuthContext = TAuthContext, State = any, Params = any> = {
  context: TBearerLambdaContext<AuthContext>
  params: Params
  state: Partial<State>
}

/**
 * Fetch any data
 */
export type TFetchAction<AuthContext = TAuthContext, Params = Record<string, any>, ReturnedData = any> = (
  event: TFetchActionEvent<AuthContext, Params>
) => Promise<TFetchPayload<ReturnedData>>

export type TFetchActionEvent<AuthContext = TAuthContext, Params = Record<string, any>> = {
  context: TBearerLambdaContext<AuthContext>
  params: Params
}

export type TStateData = AxiosResponse<{
  Item: any
}>

export type TLambdaEvent<T = TAuthContext> = {
  queryStringParameters: Record<string, any>
  context: Record<string, any> & TBearerLambdaContext<T>
  body?: any
}

export type TLambdaCallback = (error: any | null, data: any) => void
