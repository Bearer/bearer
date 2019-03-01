import Bearer, { TBearerOptions } from './bearer'
import i18n, { I18n } from './i18n'

const bearer: TBearer = (token: string, options?: Partial<TBearerOptions>) => new Bearer(token, options)

bearer.version = 'BEARER_VERSION'
bearer.i18n = i18n

export type TBearer = {
  (token: string, options?: Partial<TBearerOptions>): Bearer
  version: string
  i18n: I18n
}

export default bearer
