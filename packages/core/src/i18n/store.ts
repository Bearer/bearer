import merge from 'lodash.merge'
import get from 'lodash.get'

const DEFAULT_LOCALE = 'en'
export interface I18nStore {
  get: (key: string) => string | undefined
  setLocale: (locale: string) => void
  loadLocale: (locale: string, dictionnary: TransLationObject) => void
}

export class Store implements I18nStore {
  constructor(private readonly translationStore = window.bearerI18nStore || {}) {
    this.refreshLocale()
  }

  get(key: string): string {
    return get(this.translationStore, [this.locale, key].join('.'))
  }

  setLocale(locale: string): void {
    this.locale = locale
  }

  loadLocale(locale: string, dictionnary: TransLationObject): void {
    this.translationStore[locale] = merge(get(locale, this.translationStore), dictionnary)
  }

  refreshLocale = () => {
    this.locale = this.detectLocale()
  }

  private locale: string
  private detectLocale(): string {
    if (!window) {
      return DEFAULT_LOCALE
    }

    const browserLamg = navigator.languages
      .map(lang => lang.split('-')[0].toLowerCase())
      .find(locale => this.translationStore[locale])
    return (browserLamg || document.documentElement.lang || DEFAULT_LOCALE).toLowerCase()
  }
}

type TransLationValue = string | number | TransLationObject

interface TransLationObject {
  [key: string]: TransLationValue
}

declare const window: Window & { bearerI18nStore: any }

export default new Store()
