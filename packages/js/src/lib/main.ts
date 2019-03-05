import Bearer, { TBearerOptions } from './bearer'
import i18n, { I18n } from './i18n'

/**
 * @param  {string} clientId Client ID you'll find within the developer portal > Settings
 * @param  {Partial<TBearerOptions>} options? Fine tune bearer functionalities
 */
const bearer: TBearer = (clientId: string, options?: Partial<TBearerOptions>) => {
  bearer.instance = new Bearer(clientId, options)
  return bearer.instance
}

bearer.version = 'BEARER_VERSION'
bearer.i18n = i18n

export type TBearer = {
  (clientId: string, options?: Partial<TBearerOptions>): Bearer
  instance?: Bearer
  version: string
  i18n: I18n
}

export default bearer
