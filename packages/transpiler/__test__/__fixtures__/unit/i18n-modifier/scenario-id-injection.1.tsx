import { Component, t, p } from '@bearer/core'

@Component({
  tag: 'do-nothing'
})
class ComponentRequiringAliasing {
  fromProperty = () => {
    return (
      <span>
        <bearer-i18n key="something.key" default="MyDefault string" />
      </span>
    )
  }

  get myTranslatedTitle(): string {
    return t('my.key', 'and my default', { ok: 'ko' })
  }

  render() {
    return (
      <div>
        {t('render.t', 'translated')}
        {p('render.p', 'pluralized')}
        <bearer-i18n key="other.key" default="from the render" />
      </div>
    )
  }
}
