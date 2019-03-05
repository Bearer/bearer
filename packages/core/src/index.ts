import { Bearer as TBearer } from '@bearer/js'

import Bearer from './bearer'
import EventNames from './event-names'
import * as Requests from './requests'
import { TTranslatorFunc, TPluralizerFunc, scopedPluralize, scopedTranslate } from './i18n/index'
export { TTranslatorFunc, TPluralizerFunc } from './i18n/index'

import * as bearerState from './bearer-state'
export * from './decorators'
export const StateManager = bearerState
export const Events = EventNames
export const requests = Requests

// Next 2 helpers get rewritten with transpiler

export declare const t: TTranslatorFunc

export declare const p: TPluralizerFunc

declare const window: Window & { bearer: TBearer }

export const scopedT = (scope: string) => scopedTranslate(scope)(window.bearer.i18n)
export const scopedP = (scope: string) => scopedPluralize(scope)(window.bearer.i18n)

export default Bearer

if (process.env.BUILD !== 'distribution') {
  console.warn(`[BEARER] Running non production Bearer Core lib | version ${Bearer.version}`)
}
