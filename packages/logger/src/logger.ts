import debug from 'debug'

export default (scope: string) => debug(`bearer:${scope}`)
