import { Component, t as tAliased, p as pAliased } from '@bearer/core'

@Component({
  tag: 'do-nothing'
})
class ComponentWithAliasedImports {
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
        {tAliased('render.t', 'translated')}
        {pAliased('render.p', 'pluralized')}
        <bearer-i18n key="other.key" default="from the render" />
      </div>
    )
  }
}
