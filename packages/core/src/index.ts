import Bearer from './bearer'
import classNames from './classnames'
import * as Debug from './debug'
import EventNames from './event-names'
import * as Requests from './requests'
import { translate, pluralize, scopedPluralize, scopedTranslate } from './i18n/index'

import * as bearerState from './bearer-state'
export * from './decorators'
export const StateManager = bearerState
export const Events = EventNames
export const requests = Requests
export const debug = Debug
export const classnames = classNames

// Next 2 helpers get rewritten with transpiler
/**
 * t: i18n helper function that let you translate text easily
 * @param {string} key - Key to use for translation ex: titles.welcome.
 * @param {string} defaultValue -  A default value to use until the key get tranlated
 * @param {object} vars - An object with all required keys to replace from the template.
 */
export const t = translate
/**
 * p: i18n helper function that let you pluralize text easily
 * @param {string} key - Key to use for translation ex: titles.welcome.
 * @param {number} count - Value used as discriminator for translation.
 * @param {string} defaultValue - A default value to use until the key get tranlated.
 * @param {object} vars - An object with all required keys to replace from the template.
 */
export const p = pluralize

export const scopedT = (scope: string) => scopedTranslate(scope)(Bearer.i18nStore)
export const scopedP = (scope: string) => scopedPluralize(scope)(Bearer.i18nStore)

export default Bearer

if (process.env.BUILD !== 'distribution') {
  console.warn(`[BEARER] Running non production Bearer Core lib | version ${Bearer.version}`)
}
