import Bearer from './Bearer'
import classNames from './classnames'
import * as Debug from './debug'
import EventNames from './EventNames'
import * as Requests from './requests'

import * as bearerState from './BearerState'
export * from './decorators'
export const StateManager = bearerState
export const Events = EventNames
export const requests = Requests
export const debug = Debug
export const classnames = classNames

export default Bearer

if (process.env.BUILD !== 'distribution') {
  console.warn(`[BEARER] Running non production Bearer Core lib | version ${Bearer.version}`)
}
