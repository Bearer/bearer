import { Bearer as TBearer } from '@bearer/js'

import Bearer from './bearer'
import EventNames from './event-names'
import { TTranslatorFunc, TPluralizerFunc, scopedPluralize, scopedTranslate } from './i18n/index'
export { TTranslatorFunc, TPluralizerFunc } from './i18n/index'

export * from './decorators'
export const Events = EventNames

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
