import { BearerFetch, BearerState, Intent, IntentType, RootComponent } from '@bearer/core'

@RootComponent({
  group: 'attach-pull-request',
  role: 'display',
  shadow: false
})
export class InvalidRootComponent {
  @Intent('ListRepositories')
  fetcher: BearerFetch
  @BearerState()
  attachedPullRequests: Array<any>
  @BearerState({ statePropName: 'goats' })
  ducks: Array<any>
  @Intent('getPullRequest', IntentType.FetchData)
  fetchResource: BearerFetch

  render() {
    return <bearer-scrollable fetcher={this.fetcher} />
  }
}
