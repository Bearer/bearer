export enum Authentications {
  OAuth2 = 'OAUTH2',
  OAuth1 = 'OAUTH1',
  Basic = 'BASIC',
  ApiKey = 'APIKEY',
  NoAuth = 'NONE',
  Custom = 'CUSTOM'
}

export type TAuthentications = Record<Authentications, { name: string; value: Authentications }>

export default Authentications
