import { Component } from '@bearer/core'

@Component({
  tag: '{{componentTagName}}',
  styleUrl: '{{viewName}}.css',
  shadow: true
})
export class {{viewName}} {
  render() {
    return (
      <div class="root">
        <bearer-typography as="h2">
          {{viewName}} Component
        </bearer-typography>
      </div>
    )
  }
}
