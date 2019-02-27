import Bearer, { TBearerOptions } from './bearer'
import i18n from './i18n'

const bearer = (token: string, options?: Partial<TBearerOptions>) => new Bearer(token, options)

bearer.version = 'BEARER_VERSION'
bearer.i18n = i18n

export default bearer
