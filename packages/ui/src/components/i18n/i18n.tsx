import { Component, t, p, Prop, Events, State, Watch } from '@bearer/core'

@Component({
  tag: 'bearer-i18n'
})
export class BearerI18n {
  @Prop() _key: string // react hack
  @Prop() key: string
  @Prop() default?: string
  @Prop() var?: any
  @Prop() count?: number

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

  componentWillLoad() {
    this.dataDidChangeHandler(this.var)
  }

  componentDidLoad() {
    document.addEventListener(Events.LOCALE_CHANGED, (e: CustomEvent) => {
      this.locale = e.detail.locale
    })
  }

  componentDidUnload() {
    document.removeEventListener(Events.LOCALE_CHANGED, (this as any).forceUpdate)
  }

  render() {
    const key = this.key || this._key
    if (isNaN(this.count)) {
      return t(key, this.default, this.innerVar)
    }
    return p(key, this.count, this.default, this.innerVar)
  }
}
