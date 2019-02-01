import merge from 'lodash.merge'
import get from 'lodash.get'
import Events from '../event-names'

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
    this.localeChanged()
  }

  loadLocale(locale: string, dictionnary: TransLationObject): void {
    this.translationStore[locale] = merge(get(locale, this.translationStore), dictionnary)
    this.localeChanged()
  }

  refreshLocale = () => {
    this.locale = this.detectLocale()
  }

  private localeChanged() {
    document.dispatchEvent(new CustomEvent(Events.LOCALE_CHANGED, { detail: { locale: this.locale } }))
  }

  private locale: string
  private detectLocale(): string {
    if (!window) {
      return DEFAULT_LOCALE
    }

    const browserLang = navigator.languages
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
