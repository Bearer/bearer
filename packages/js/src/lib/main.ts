import Bearer from './bearer'

const bearer = (token: string) => new Bearer(token)

bearer.version = 'BEARER_VERSION'

export default bearer
