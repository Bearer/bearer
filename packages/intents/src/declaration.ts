import { AxiosResponse } from 'axios'

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

export type TStateData = AxiosResponse<{
  Item: any
}>
