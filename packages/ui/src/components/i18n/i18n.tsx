import { Component, t, p, Prop, Events, State, Watch } from '@bearer/core'

@Component({
  tag: 'bearer-i18n'
})
export class BearerI18n {
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
    if (isNaN(this.count)) {
      return t(this.key, this.default, this.innerVar)
    }
    return p(this.key, this.count, this.default, this.innerVar)
  }
}
