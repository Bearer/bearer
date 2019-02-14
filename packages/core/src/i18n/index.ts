import { I18nStore } from './store'
import template from 'lodash.template'

const interpolate = /{{([\s\S]+?)}}/g

export const translate = (store: I18nStore) => (key: string, defaultValue: string, vars?: Record<string, any>) => {
  return template(store.get(key) || defaultValue, {
    interpolate
  })(vars || {})
}

export const pluralize = (store: I18nStore) => (
  key: string,
  count: number,
  defaultValue: string,
  vars?: Record<string, any>
) => {
  const keyWithCount = [key, count].join('.')
  if (store.get(keyWithCount)) {
    return translate(store)(keyWithCount, defaultValue, vars)
  }
  const quantity = count > 1 ? 'many' : count
  const newKey = [key, quantity].join('.')

  return translate(store)(newKey, defaultValue, { count, ...vars })
}

export default {
  translate,
  pluralize
}
