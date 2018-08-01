import { RootComponent, Intent, BearerFetch, IntentType, BearerState } from '@bearer/core'

@RootComponent({
  group: 'AttachPullRequest',
  name: 'action'
})
export class AttachPullRequestAction {
  @Intent('ListRepositories') fetcher: BearerFetch
  @RetrieveStateIntent() retrieve: BearerFetch
  @BearerState() attachedPullRequests: Array<any>
  @BearerState({ statePropName: 'goats' })
  ducks: Array<any>
  @Intent('getPullRequest', IntentType.GetResource)
  fetchResource: BearerFetch

  render() {
    return <bearer-scrollable fetcher={this.fetcher} />
  }
}
