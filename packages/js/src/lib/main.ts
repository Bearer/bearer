import Bearer, { TBearerOptions } from './bearer'
import i18n, { I18n } from './i18n'

const bearer: TBearer = (token: string, options?: Partial<TBearerOptions>) => {
  bearer.instance = new Bearer(token, options)
  return bearer.instance
}

bearer.version = 'BEARER_VERSION'
bearer.i18n = i18n

export type TBearer = {
  (token: string, options?: Partial<TBearerOptions>): Bearer
  instance?: Bearer
  version: string
  i18n: I18n
}

export default bearer
