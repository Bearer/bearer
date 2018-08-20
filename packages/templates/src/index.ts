import Authentication from "@bearer/types/lib/Authentications";

import apikey from './apikey'
import basicauth from './basicauth'
import noauth from './noauth'
import oauth2 from './oauth2'

type  TTemplate = {
  SaveState: string,
  RetrieveState: string,
  FetchData: string
}

export const templates: Record<Authentication, TTemplate> = {
  [Authentication.OAuth2]: oauth2,
  [Authentication.ApiKey]: apikey,
  [Authentication.Basic]: basicauth,
  [Authentication.NoAuth]: noauth,
}

export default templates
