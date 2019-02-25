import Bearer from './bearer'
import i18n from './i18n'

const bearer = (token: string) => new Bearer(token)

bearer.version = 'BEARER_VERSION'
bearer.i18n = i18n

export default bearer
