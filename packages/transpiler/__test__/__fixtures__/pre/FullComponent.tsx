import { Component, Intent, BearerFetch, IntentType } from '@bearer/core'

@Component({
  tag: 'list-repositories',
  styleUrl: 'ListRepositories.css',
  shadow: true
})
export class ListRepositories {
  @Intent('ListRepositories') fetcher: BearerFetch

  @Intent('getPullRequest', IntentType.GetResource)
  fetchResource: BearerFetch

  render() {
    return <bearer-scrollable fetcher={this.fetcher} />
  }
}
