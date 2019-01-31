import { I18nStore } from './store'

export default (store: I18nStore) => (key: string, defaultValue: string) => {
  return store.get(key) || defaultValue
}
