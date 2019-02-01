import { I18nStore } from './store'
import template from 'lodash.template'

const interpolate = /{{([\s\S]+?)}}/g

export default (store: I18nStore) => (key: string, defaultValue: string, vars?: Record<string, any>) => {
  return template(store.get(key) || defaultValue, {
    interpolate
  })(vars || {})
}
