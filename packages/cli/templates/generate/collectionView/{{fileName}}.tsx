import { Component, Intent, BearerFetch } from '@bearer/core'

@Component({
  tag: '{{componentTagName}}',
  styleUrl: '{{componentFileName}}.css',
  shadow: true
})
export class {{componentName}} {
  @Intent('{{componentName}}') fetcher: BearerFetch
  render() {
    return <bearer-scrollable fetcher={this.fetcher} />
  }
}