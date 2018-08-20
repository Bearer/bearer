import { BearerFetch, Component, Intent } from '@bearer/core'

@Component({
  tag: '{{componentTagName}}',
  shadow: true
})
export class {{componentClassName}} {
  @Intent('{{fileName}}') fetcher: BearerFetch
  render() {
    return <bearer-scrollable fetcher={this.fetcher} />
  }
}