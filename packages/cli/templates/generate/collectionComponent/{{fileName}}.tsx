import { BearerFetch, Component, Function } from '@bearer/core'

@Component({
  tag: '{{componentTagName}}',
  shadow: true
})
export class {{componentClassName}} {
  @Function('{{fileName}}') fetcher: BearerFetch
  render() {
    return <bearer-scrollable fetcher={this.fetcher} />
  }
}
