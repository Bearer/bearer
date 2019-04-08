import { AxiosResponse } from 'axios'
import { Store } from './store'
/**
 * Simplest error payload, better add code and message so that debugging is made easier
 */
export type TErrorPayload<ReturnedError> = { error: ReturnedError }
export type TDataPayload<ReturnedData> = { data: ReturnedData; referenceId?: string }

export type TFetchPayload<ReturnedData = any, ReturnedError = any> = Partial<TDataPayload<ReturnedData>> &
  Partial<TErrorPayload<ReturnedError>>

/**
 * Contexts
 */

type TBaseAuthContext<TAuthAccessContent> = { authAccess: TAuthAccessContent; [key: string]: any }

export namespace contexts {
  export interface OAuth2 {
    accessToken: string
  }

  export interface OAuth1 {
    accessToken: string
    tokenSecret: string
  }

  export interface Basic {
    username: string
    password: string
  }

  export interface ApiKey {
    apiKey: string
  }

  export interface Custom {}
  export interface None {}
}

export type TOAUTH1AuthContext = TBaseAuthContext<contexts.OAuth1> & {
  setup: { consumerKey: string; consumerSecret: string }
}

export type TOAUTH2AuthContext = TBaseAuthContext<contexts.OAuth2>
export type TNONEAuthContext = TBaseAuthContext<contexts.None>
export type TCUSTOMAuthContext = TBaseAuthContext<contexts.Custom>
export type TBASICAuthContext = TBaseAuthContext<contexts.Basic>
export type TAPIKEYAuthContext = TBaseAuthContext<contexts.ApiKey>

export type TAuthContext =
  | TNONEAuthContext
  | TCUSTOMAuthContext
  | TOAUTH2AuthContext
  | TBASICAuthContext
  | TAPIKEYAuthContext
  | TOAUTH1AuthContext

export type TBearerLambdaContext<AuthContext = TAuthContext, DataContext = {}> = DataContext &
  AuthContext & {
    bearerBaseURL: string
  } & { isBackend: boolean; signature: string }

/**
 * Functions
 */

/**
 * Fetch any data
 */
export type TFetchAction<ReturnedData = any> = (event: TFetchActionEvent) => TFetchPromise<ReturnedData>
export type TFetchPromise<ReturnedData, ReturnedError = any> = Promise<TFetchPayload<ReturnedData, ReturnedError>>

export type TFetchActionEvent<Params = any, AuthContext = TAuthContext, DataContext = {}> = {
  context: TBearerLambdaContext<AuthContext, DataContext>
  params: Params
  store: Store
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
