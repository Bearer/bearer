import bearer from '@bearer/js'

import Bearer from './bearer'
import classNames from './classnames'
import * as Debug from './debug'
import EventNames from './event-names'
import * as Requests from './requests'
import { TTranslatorFunc, TPluralizerFunc, scopedPluralize, scopedTranslate } from './i18n/index'
export { TTranslatorFunc, TPluralizerFunc } from './i18n/index'

import * as bearerState from './bearer-state'
export * from './decorators'
export const StateManager = bearerState
export const Events = EventNames
export const requests = Requests
export const debug = Debug
export const classnames = classNames

// Next 2 helpers get rewritten with transpiler

export declare const t: TTranslatorFunc

export declare const p: TPluralizerFunc

export const scopedT = (scope: string) => scopedTranslate(scope)(bearer.i18n)
export const scopedP = (scope: string) => scopedPluralize(scope)(bearer.i18n)

export default Bearer

if (process.env.BUILD !== 'distribution') {
  console.warn(`[BEARER] Running non production Bearer Core lib | version ${Bearer.version}`)
}
