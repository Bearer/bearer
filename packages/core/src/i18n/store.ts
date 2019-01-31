export interface I18nStore {
  get: (key: string) => string | undefined
}

export class Store implements I18nStore {
  get(key: string): string {
    // that way we alway return fallback
    return null
  }
}

export default new Store()
