import * as Requests from './requests'
import * as Debug from './debug'
import Bearer from './Bearer'
import classNames from './classnames'
import EventNames from './EventNames'

// export * from './userDataClient'
export * from './decorators'
import * as bearerState from './BearerState'
export const BearerState = bearerState
export const Events = EventNames
export const requests = Requests
export const debug = Debug
export const classnames = classNames

export default Bearer

if (process.env.BUILD !== 'distribution') {
  console.warn(
    `[BEARER] Running non production Bearer Core lib | version ${
      Bearer.version
    }`
  )
}
