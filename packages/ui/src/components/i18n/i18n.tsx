import { Component, TTranslatorFunc, TPluralizerFunc, scopedP, scopedT, Prop, Events, State, Watch } from '@bearer/core'

@Component({
  tag: 'bearer-i18n'
})
export class BearerI18n {
  @Prop() _key: string // react hack

  @Prop() key: string
  @Prop() default?: string
  @Prop() var?: any
  @Prop() count?: number
  @Prop() scope?: string

  @State() translate: TTranslatorFunc
  @State() pluralize: TPluralizerFunc
  @State() innerVar: any = {}
  @State() locale: string

  @Watch('var')
  dataDidChangeHandler(newValue: string) {
    switch (typeof newValue) {
      case 'string': {
        this.innerVar = JSON.parse(newValue || '{}')
        break
      }
      case 'object': {
        this.innerVar = { ...(newValue as Object) }
        break
      }
    }
  }

  refreshStore = () => {
    this.translate = scopedT(this.scope)
    this.pluralize = scopedP(this.scope)
  }

  componentWillLoad() {
    this.dataDidChangeHandler(this.var)
    this.refreshStore()
  }

  localeUpdate = (e: CustomEvent) => {
    this.refreshStore()
    this.locale = e.detail.locale
  }

  componentDidLoad() {
    document.addEventListener(Events.LOCALE_CHANGED, this.localeUpdate)
  }

  componentDidUnload() {
    document.removeEventListener(Events.LOCALE_CHANGED, this.localeUpdate)
  }

  render() {
    const key = this.key || this._key
    if (isNaN(this.count)) {
      return this.translate(key, this.default, this.innerVar)
    }
    return this.pluralize(key, this.count, this.default, this.innerVar)
  }
}
