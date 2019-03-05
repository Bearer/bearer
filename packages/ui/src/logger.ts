import debug from '@bearer/logger'

export default (scope: string) => debug('ui').extend(scope)
