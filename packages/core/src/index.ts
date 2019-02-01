import Bearer from './bearer'
import classNames from './classnames'
import * as Debug from './debug'
import EventNames from './event-names'
import * as Requests from './requests'
import { translate, pluralize } from './i18n/index'

import * as bearerState from './bearer-state'
export * from './decorators'
export const StateManager = bearerState
export const Events = EventNames
export const requests = Requests
export const debug = Debug
export const classnames = classNames

export const t = translate(Bearer.i18nStore)
export const p = pluralize(Bearer.i18nStore)

export default Bearer

if (process.env.BUILD !== 'distribution') {
  console.warn(`[BEARER] Running non production Bearer Core lib | version ${Bearer.version}`)
}
