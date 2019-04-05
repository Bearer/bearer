export enum Authentications {
  OAuth2 = 'OAUTH2',
  OAuth1 = 'OAUTH1',
  Basic = 'BASIC',
  ApiKey = 'APIKEY',
  NoAuth = 'NONE',
  Custom = 'CUSTOM'
}

export type TAuthentications = Record<Authentications, { name: string; value: Authentications }>

// Any way to auto generate this?

export type TConfig =
  | configs.TOAuth2Config
  | configs.TOAuth1Config
  | configs.TApiKeyConfig
  | configs.TBasicConfig
  | configs.TNoAuthConfig
  | configs.TCustomConfig

export namespace configs {
  export type TBaseConfig<authType extends Authentications> = {
    authType: authType
  }

  export type TOAuth2Config = TBaseConfig<Authentications.OAuth2> & {
    authorizationURL: string
    tokenURL: string
    tokenParams: any
    authorizationParams: any
    config: {
      scope: string[]
    }
  }

  export type TOAuth1Config = TBaseConfig<Authentications.OAuth1> & {
    requestTokenURL: string
    accessTokenURL: string
    userAuthorizationURL: string
    callbackURL: string
    signatureMethod: 'HMAC-SHA1'
    tokenParams: Record<string, number | string>
    authorizationParams: Record<string, number | string>
    config: {
      scope: string[]
    }
  }

  export type TApiKeyConfig = TBaseConfig<Authentications.ApiKey> & {}
  export type TBasicConfig = TBaseConfig<Authentications.Basic> & {}
  export type TNoAuthConfig = TBaseConfig<Authentications.NoAuth> & {}
  export type TCustomConfig = TBaseConfig<Authentications.Custom> & {}
}

export default Authentications
