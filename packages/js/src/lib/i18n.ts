import merge from 'lodash.merge'
import get from 'lodash.get'

import logger from './logger'

const debug = logger.extend('i18n')
const DEFAULT_LOCALE = 'en'

export const LOCALE_CHANGED = 'bearer-locale-changed'

export class I18n {
  private _locale = DEFAULT_LOCALE
  private _dictionnary: TI18nDictionnary = {
    [DEFAULT_LOCALE]: {}
  }

  load = async (
    integrationName: string | null,
    dictionnary:
      | TransLationObject
      | IntegrationTranlsationEntry
      | Promise<TransLationObject | IntegrationTranlsationEntry>,
    { locale = this.locale }: Partial<{ locale: string }> = {}
  ) => {
    const result = await dictionnary
    const newEntries = !!integrationName
      ? { [integrationName]: result as TransLationObject }
      : (result as IntegrationTranlsationEntry)

    this._dictionnary[locale] = merge(get(this._dictionnary, locale), newEntries)
    this.localeChanged()
  }

  private localeChanged() {
    document.dispatchEvent(new CustomEvent(LOCALE_CHANGED, { detail: { locale: this.locale } }))
  }

  // @ts-ignore
  get = (integrationName: string | null, key: string, options: Partial<{ locale: string }> = {}): TransLationValue => {
    const path = [options.locale || this.locale, integrationName, key].filter(m => m).join('.')
    debug('lookup key', path)
    return get(this._dictionnary, path)
  }

  set locale(locale: string) {
    this._locale = locale
    this.localeChanged()
  }

  get locale(): string {
    return this._locale
  }
}

const i18n = new I18n()

type TI18nDictionnary = {
  [locale: string]: IntegrationTranlsationEntry
}

type IntegrationTranlsationEntry = { [integrationName: string]: TransLationObject }

type TransLationValue = string | number | TransLationObject

interface TransLationObject {
  [key: string]: TransLationValue
}

export default i18n
