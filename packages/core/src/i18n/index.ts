import { I18nStore } from './store'
import template from 'lodash.template'

const interpolate = /{{([\s\S]+?)}}/g

export const scopedTranslate = (scope?: string): TTranslator => (store: I18nStore) => (
  key: string,
  defaultValue: string,
  vars?: Record<string, any>
) => {
  return template(store.get(key, scope) || defaultValue, {
    interpolate
  })(vars || {})
}

export const scopedPluralize = (scope?: string): TPluralizer => (store: I18nStore) => (
  key: string,
  count: number,
  defaultValue: string,
  vars?: Record<string, any>
) => {
  const keyWithCount = [key, count].join('.')
  if (store.get(keyWithCount, scope)) {
    return scopedTranslate(scope)(store)(keyWithCount, defaultValue, vars)
  }
  const quantity = count > 1 ? 'many' : count
  const newKey = [key, quantity].join('.')
  return scopedTranslate(scope)(store)(newKey, defaultValue, vars)
}

export type TTranslatorFunc = {
  (key: string, defaultValue: string, vars?: Record<string, any>): string
}

type TTranslator = {
  (store: I18nStore): TTranslatorFunc
}

export type TPluralizerFunc = {
  (key: string, count: number, defaultValue: string, vars?: Record<string, any>): string
}
type TPluralizer = {
  (store: I18nStore): TPluralizerFunc
}
